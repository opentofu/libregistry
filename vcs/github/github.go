// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/opentofu/libregistry/logger"
	"github.com/opentofu/libregistry/vcs"
)

// New creates a new GitHub VCS client.
func New(
	options ...Opt,
) (vcs.Client, error) {
	config := Config{}
	for _, opt := range options {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}
	config.ApplyDefaults()

	return &github{
		config: config,
		lock:   &sync.Mutex{},
		locks:  map[string]*sync.Mutex{},
	}, nil
}

type github struct {
	config Config
	lock   *sync.Mutex
	locks  map[string]*sync.Mutex
}

func (g github) GetTagVersion(ctx context.Context, repository vcs.RepositoryAddr, version vcs.VersionNumber) (vcs.Version, error) {
	if err := repository.Validate(); err != nil {
		return vcs.Version{}, err
	}
	if err := version.Validate(); err != nil {
		return vcs.Version{}, err
	}

	wc, err := g.getWorkingCopy(ctx, repository)
	if err != nil {
		return vcs.Version{}, nil
	}
	defer wc.cleanup()
	return wc.getTag(ctx, version)
}

func (g github) GetRepositoryBrowseURL(_ context.Context, repository vcs.RepositoryAddr) (string, error) {
	if err := repository.Validate(); err != nil {
		return "", err
	}
	return "https://github.com/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name), nil
}

func (g github) GetVersionBrowseURL(_ context.Context, repository vcs.RepositoryAddr, version vcs.VersionNumber) (string, error) {
	if err := repository.Validate(); err != nil {
		return "", err
	}
	if err := version.Validate(); err != nil {
		return "", err
	}
	return "https://github.com/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name) + "/tree/" + url.PathEscape(string(version)), nil
}

func (g github) GetFileViewURL(_ context.Context, repository vcs.RepositoryAddr, version vcs.VersionNumber, file string) (string, error) {
	if err := repository.Validate(); err != nil {
		return "", err
	}
	if err := version.Validate(); err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("empty file name passed")
	}
	fileParts := strings.Split(file, "/")
	for i, part := range fileParts {
		fileParts[i] = url.PathEscape(part)
	}
	return "https://github.com/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name) + "/blob/" + url.PathEscape(string(version)) + "/" + file, nil
}

func (g github) GetRepositoryInfo(ctx context.Context, repository vcs.RepositoryAddr) (vcs.RepositoryInfo, error) {
	if err := repository.Validate(); err != nil {
		return vcs.RepositoryInfo{}, err
	}
	type repoInfoResponse struct {
		Description string `json:"description"`
	}

	var response repoInfoResponse

	if err := g.request(ctx, "https://api.github.com/repos/"+url.PathEscape(string(repository.Org))+"/"+url.PathEscape(repository.Name), &response); err != nil {
		return vcs.RepositoryInfo{}, err
	}
	return vcs.RepositoryInfo{
		Description: response.Description,
	}, nil
}

func (g github) Checkout(ctx context.Context, repository vcs.RepositoryAddr, version vcs.VersionNumber) (vcs.WorkingCopy, error) {
	if err := version.Validate(); err != nil {
		return nil, err
	}
	wc, err := g.getWorkingCopy(ctx, repository)
	if err != nil {
		return nil, err
	}

	tagExists, err := wc.tagExists(ctx, version)
	if err != nil {
		wc.cleanup()
		return nil, err
	}
	if !tagExists {
		wc.cleanup()
		return nil, &vcs.VersionNotFoundError{
			RepositoryAddr: repository,
			Version:        version,
		}
	}

	if err := wc.checkout(ctx, version); err != nil {
		wc.cleanup()
		return nil, err
	}
	return wc, nil
}

func (g github) getWorkingCopy(ctx context.Context, repository vcs.RepositoryAddr) (*workingCopy, error) {
	if err := repository.Validate(); err != nil {
		return nil, err
	}
	parentDirectory := path.Join(g.config.CheckoutRootDirectory, string(repository.Org))
	checkoutDirectory := path.Join(parentDirectory, repository.Name)
	gitDirectory := path.Join(checkoutDirectory, ".git")

	g.lock.Lock()
	lock, ok := g.locks[checkoutDirectory]
	if !ok {
		lock = &sync.Mutex{}
		g.locks[checkoutDirectory] = lock
	}
	g.lock.Unlock()
	lock.Lock()
	cleanup := func() {
		g.lock.Lock()

		if !g.config.SkipCleanupWorkingCopyOnClose {
			// Make sure that any open file descriptors are closed before cleaning up the directory so Windows file
			// locking doesn't block the cleanup:
			runtime.GC()

			if err := os.RemoveAll(checkoutDirectory); err != nil {
				g.config.Logger.Debug(ctx, "Failed to clean up clone repository at %s (%v)", checkoutDirectory, err)
			}
		}

		delete(g.locks, checkoutDirectory)
		lock.Unlock()
		g.lock.Unlock()
	}

	stat, err := os.Stat(gitDirectory)
	if err != nil || !stat.IsDir() {
		if err := os.RemoveAll(checkoutDirectory); err != nil {
			cleanup()
			return nil, fmt.Errorf("failed to remove broken checkout directory %s (%w)", checkoutDirectory, err)
		}
		if err := os.MkdirAll(parentDirectory, 0700); err != nil {
			cleanup()
			return nil, fmt.Errorf("failed to create checkout parent directory %s (%w)", parentDirectory, err)
		}
		credentials := ""
		if g.config.Username != "" && g.config.Token != "" {
			credentials = url.PathEscape(g.config.Username) + ":" + url.PathEscape(g.config.Token) + "@"
		}

		cloneURL := "https://" + credentials + "github.com/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name) + ".git"
		if err := g.git(ctx, parentDirectory, nil, "clone", "--depth", "1", cloneURL, checkoutDirectory); err != nil {
			cleanup()

			// Clone failed, check if repository exists.
			repoExists, e := g.repositoryExists(ctx, repository)
			if e == nil && !repoExists {
				return nil, &vcs.RepositoryNotFoundError{RepositoryAddr: repository, Cause: err}
			}

			return nil, err
		}
	}

	if err := g.git(ctx, checkoutDirectory, nil, "fetch", "--tags"); err != nil {
		cleanup()
		return nil, err
	}

	return &workingCopy{
		ReadDirFS:  os.DirFS(checkoutDirectory).(fs.ReadDirFS),
		repository: repository,
		dir:        checkoutDirectory,
		cleanup:    cleanup,
		g:          g,
	}, nil
}

type workingCopy struct {
	fs.ReadDirFS
	cleanup    func()
	repository vcs.RepositoryAddr
	dir        string
	g          github
}

func (w workingCopy) RawDirectory() (string, error) {
	return w.dir, nil
}

func (w workingCopy) Close() error {
	w.cleanup()
	return nil
}

func (w workingCopy) checkout(ctx context.Context, version vcs.VersionNumber) error {
	if err := w.g.git(ctx, w.dir, nil, "checkout", string(version)); err != nil {
		// Checkout failed, see if tag exists.
		tagExists, e := w.tagExists(ctx, version)
		if e == nil && !tagExists {
			return &vcs.VersionNotFoundError{Version: version, RepositoryAddr: w.repository, Cause: err}
		}

		return err
	}
	return nil
}

func (w workingCopy) tagExists(ctx context.Context, version vcs.VersionNumber) (bool, error) {
	tags, err := w.listTags(ctx)
	if err != nil {
		return false, err
	}
	for _, tag := range tags {
		if tag.VersionNumber.Equals(version) {
			return true, nil
		}
	}
	return false, nil
}

func (w workingCopy) getTag(ctx context.Context, tag vcs.VersionNumber) (vcs.Version, error) {
	tags, err := w.listTags(ctx)
	if err != nil {
		return vcs.Version{}, err
	}
	for _, t := range tags {
		if t.VersionNumber.Equals(tag) {
			return t, nil
		}
	}
	return vcs.Version{}, vcs.VersionNotFoundError{
		RepositoryAddr: w.repository,
		Version:        tag,
	}
}

func (w workingCopy) listTags(ctx context.Context) ([]vcs.Version, error) {
	stdout := &bytes.Buffer{}
	if err := w.g.git(ctx, w.dir, stdout, "for-each-ref", "--format=%(refname:short)\t%(creatordate:format:%s)", "refs/tags/*"); err != nil {
		return nil, err
	}
	lines := strings.Split(stdout.String(), "\n")
	var result []vcs.Version
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line does not contain enough parts to parse: %s", line)
		}
		tag := vcs.VersionNumber(strings.ReplaceAll(parts[0], "refs/tags/", ""))
		unixTime, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse git output: %s (%v)", line, err)
		}
		created := time.Unix(int64(unixTime), 0)
		ver := vcs.Version{
			VersionNumber: tag,
			Created:       created,
		}
		if err := ver.Validate(); err != nil {
			w.g.config.Logger.Debug(ctx, "Skipping tag %s because it does not match the naming rules.", ver.VersionNumber)
			continue
		}
		result = append(result, ver)
	}
	return result, nil
}

func (g github) git(ctx context.Context, dir string, stdout io.Writer, params ...string) error {
	cmd := exec.Command(g.config.GitPath, params...)
	commandString := strings.Join(append([]string{g.config.GitPath}, params...), " ")
	logger.LogTrace(ctx, g.config.Logger, "Running "+commandString)
	if stdout == nil {
		stdout = logger.NewWriter(ctx, g.config.Logger, logger.LevelDebug, commandString+": ")
	}
	cmd.Stdout = stdout
	cmd.Stderr = logger.NewWriter(ctx, g.config.Logger, logger.LevelDebug, commandString+": ")
	cmd.Dir = dir
	done := make(chan struct{})
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run %s (%w)", commandString, err)
	}
	var lastErr error
	go func() {
		defer close(done)
		lastErr = cmd.Wait()
	}()
	select {
	case <-done:
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		<-done
	}
	if lastErr == nil {
		return nil
	}
	var exitErr *exec.ExitError
	if !errors.As(lastErr, &exitErr) {
		return fmt.Errorf("%s failed (%w)", commandString, lastErr)
	}
	if exitErr.ExitCode() != 0 {
		return fmt.Errorf("%s exited with exit code %d", commandString, exitErr.ExitCode())
	}
	return nil
}

func (g github) ParseRepositoryAddr(ref string) (vcs.RepositoryAddr, error) {
	ref = strings.TrimPrefix(ref, "github.com/")
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) != 2 {
		return vcs.RepositoryAddr{}, &vcs.InvalidRepositoryAddrError{
			RepositoryString: ref,
		}
	}
	result := vcs.RepositoryAddr{
		Org:  vcs.OrganizationAddr(parts[0]),
		Name: parts[1],
	}
	return result, result.Validate()
}

type rss struct {
	Entry []struct {
		ID      string `title:"id"`
		Title   string `xml:"title"`
		Updated string `xml:"updated"`
	} `xml:"entry"`
}

func (g github) ListLatestTags(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	logger.LogTrace(ctx, g.config.Logger, "Requesting latest tags for repository %s...", repository)
	return g.listLatest(ctx, repository, "tags.atom")
}

func (g github) ListLatestReleases(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	logger.LogTrace(ctx, g.config.Logger, "Requesting latest releases for repository %s...", repository)
	return g.listLatest(ctx, repository, "releases.atom")
}

func (g github) listLatest(ctx context.Context, repository vcs.RepositoryAddr, file string) ([]vcs.Version, error) {
	if err := repository.Validate(); err != nil {
		return nil, err
	}
	rssURL := "https://github.com/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name) + "/" + file
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rssURL, nil)
	if err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: fmt.Errorf("invalid HTTP request (%w)", err),
		}
	}
	resp, err := g.config.HTTPClient.Do(req)
	if err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: err,
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, &vcs.RequestFailedError{
			Cause: fmt.Errorf("invalid status code: %d", resp.StatusCode),
			Body:  body,
		}
	}

	decoder := xml.NewDecoder(resp.Body)
	response := rss{}
	if err := decoder.Decode(&response); err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: fmt.Errorf("failed to decode RSS (%w)", err),
		}
	}

	var result []vcs.Version
	for _, entry := range response.Entry {
		versionNumber := vcs.VersionNumber(entry.Title)
		if err = versionNumber.Validate(); err != nil {
			g.config.Logger.Debug(ctx, "Skipping invalid version %s when querying %s in repository %s", versionNumber, file, repository)
			continue
		}
		versionCreated, err := time.Parse(time.RFC3339, entry.Updated)
		if err != nil {
			g.config.Logger.Debug(ctx, "Skipping invalid creation version creation time %s when querying %s in repository %s", entry.Updated, file, repository)
			continue
		}

		result = append(result, vcs.Version{
			VersionNumber: versionNumber,
			Created:       versionCreated,
		})
	}
	return result, nil
}

func (g github) ListAllTags(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	logger.LogTrace(ctx, g.config.Logger, "Requesting all tags for repository %s...", repository)

	wc, err := g.getWorkingCopy(ctx, repository)
	if err != nil {
		return nil, err
	}
	defer wc.cleanup()
	return wc.listTags(ctx)
}

func (g github) ListAllReleases(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	logger.LogTrace(ctx, g.config.Logger, "Requesting all releases for repository %s...", repository)
	return g.listAll(ctx, repository, "release")
}

func (g github) listAll(ctx context.Context, repository vcs.RepositoryAddr, itemType string) ([]vcs.Version, error) {
	type responseItem struct {
		Name vcs.VersionNumber `json:"name"`
		// PublishedAt is only present for releases, not for tags.
		PublishedAt string `json:"published_at,omitempty"`
	}

	if err := repository.Validate(); err != nil {
		return nil, err
	}
	reqURL := "https://api.github.com/repos/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name) + "/" + itemType + "s"

	var response []responseItem
	if err := g.request(ctx, reqURL, &response); err != nil {
		var statusCodeErr *InvalidStatusCodeError
		if errors.As(err, &statusCodeErr) && statusCodeErr.StatusCode == http.StatusNotFound {
			return nil, &vcs.RepositoryNotFoundError{
				RepositoryAddr: repository,
				Cause:          err,
			}
		}
		return nil, err
	}

	var result []vcs.Version
	for _, item := range response {
		err := item.Name.Validate()
		if err != nil {
			g.config.Logger.Debug(ctx, "Skipping invalid %s %s in repository %s", itemType, item.Name, repository)
			continue
		}
		created := time.Time{}
		if item.PublishedAt != "" {
			created, err = time.Parse(time.RFC3339, item.PublishedAt)
			if err != nil {
				g.config.Logger.Debug(ctx, "Skipping invalid %s creation date (%s) for %s in repository %s", itemType, item.PublishedAt, item.Name, repository)
				continue
			}
		}
		result = append(result, vcs.Version{
			VersionNumber: item.Name,
			Created:       created,
		})
	}
	return result, nil
}

func (g github) request(ctx context.Context, url string, response any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return &vcs.RequestFailedError{
			Cause: fmt.Errorf("invalid HTTP request (%w)", err),
		}
	}
	if g.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.config.Token)
	}
	logger.LogTrace(ctx, g.config.Logger, "Sending GET request to %s...", url)
	resp, err := g.config.HTTPClient.Do(req)
	if err != nil {
		logger.LogTrace(ctx, g.config.Logger, "GET request to %s failed (%v)", url, err)
		return &vcs.RequestFailedError{
			Cause: err,
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	logger.LogTrace(ctx, g.config.Logger, "GET request to %s returned status code %d", url, resp.StatusCode)
	if resp.StatusCode != 200 {
		return &vcs.RequestFailedError{
			Cause: &InvalidStatusCodeError{resp.StatusCode},
		}
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&response); err != nil {
		g.config.Logger.Warn(ctx, "GitHub returned an invalid JSON when requesting %s (%v)", url, err)
		return &vcs.RequestFailedError{
			Cause: fmt.Errorf("failed to decode response (%w)", err),
		}
	}

	return nil
}

func (g github) ListAssets(ctx context.Context, repository vcs.RepositoryAddr, version vcs.VersionNumber) ([]vcs.AssetName, error) {
	logger.LogTrace(ctx, g.config.Logger, "Listing assets for repository %s version %s", repository, version)

	type responseItem struct {
		Name   string `json:"name"`
		Assets []struct {
			Name vcs.AssetName `json:"name"`
		} `json:"assets"`
	}

	if err := repository.Validate(); err != nil {
		return nil, err
	}
	if err := version.Validate(); err != nil {
		return nil, err
	}

	reqURL := "https://api.github.com/repos/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name) + "/releases/tags/" + url.PathEscape(string(version))

	var response responseItem
	if err := g.request(ctx, reqURL, &response); err != nil {
		var statusCodeErr *InvalidStatusCodeError
		if errors.As(err, &statusCodeErr) && statusCodeErr.StatusCode == http.StatusNotFound {
			return nil, &vcs.VersionNotFoundError{
				RepositoryAddr: repository,
				Version:        version,
				Cause:          err,
			}
		}

		return nil, err
	}
	var result []vcs.AssetName
	for _, asset := range response.Assets {
		err := asset.Name.Validate()
		if err != nil {
			g.config.Logger.Debug(ctx, "Skipping invalid asset named %s in repository %s release %s", asset.Name, repository, version)
		} else {
			result = append(result, asset.Name)
		}
	}
	return result, nil
}

func (g github) DownloadAsset(ctx context.Context, repository vcs.RepositoryAddr, version vcs.VersionNumber, asset vcs.AssetName) ([]byte, error) {
	if err := repository.Validate(); err != nil {
		return nil, err
	}
	if err := version.Validate(); err != nil {
		return nil, err
	}
	if err := asset.Validate(); err != nil {
		return nil, err
	}
	logger.LogTrace(ctx, g.config.Logger, "Listing asset %s for repository %s version %s", asset, repository, version)
	assetURL := "https://api.github.com/repos/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name) + "/releases/download/" + url.PathEscape(string(version)) + "/" + url.PathEscape(string(asset))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, assetURL, nil)
	if err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: fmt.Errorf("invalid HTTP request (%w)", err),
		}
	}
	if g.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.config.Token)
	}
	logger.LogTrace(ctx, g.config.Logger, "Sending GET request to %s...", assetURL)
	resp, err := g.config.HTTPClient.Do(req)
	if err != nil {
		logger.LogTrace(ctx, g.config.Logger, "GET request to %s failed (%v)", assetURL, err)
		return nil, &vcs.RequestFailedError{
			Cause: err,
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	logger.LogTrace(ctx, g.config.Logger, "GET request to %s returned status code %d", assetURL, resp.StatusCode)
	if resp.StatusCode != 200 {
		err = InvalidStatusCodeError{resp.StatusCode}
		if resp.StatusCode == 404 {
			return nil, &vcs.AssetNotFoundError{
				RepositoryAddr: repository,
				Version:        version,
				Asset:          asset,
				Cause:          err,
			}
		}
		return nil, &vcs.RequestFailedError{
			Cause: err,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: err,
		}
	}
	return body, nil
}

func (g github) HasPermission(ctx context.Context, username vcs.Username, organization vcs.OrganizationAddr) (bool, error) {
	type memberType struct {
		Login string `json:"login"`
	}

	if err := organization.Validate(); err != nil {
		return false, err
	}
	if err := username.Validate(); err != nil {
		return false, err
	}
	logger.LogTrace(ctx, g.config.Logger, "Checking if user %s has permissions for the organization %s...", username, organization)
	reqURL := "https://api.github.com/orgs/" + url.PathEscape(string(organization)) + "/members"
	var response []memberType
	if err := g.request(ctx, reqURL, &response); err != nil {
		var statusCodeErr *InvalidStatusCodeError
		if errors.As(err, &statusCodeErr) && statusCodeErr.StatusCode == http.StatusNotFound {
			return false, &vcs.OrganizationNotFoundError{
				OrganizationAddr: organization,
				Cause:            err,
			}
		}
		return false, err
	}
	for _, member := range response {
		if strings.EqualFold(member.Login, string(username)) {
			return true, nil
		}
	}
	return false, nil
}

func (g github) repositoryExists(ctx context.Context, repositoryAddr vcs.RepositoryAddr) (bool, error) {
	var repoResponse any
	if err := g.request(ctx, "https://github.com/repos/"+url.PathEscape(string(repositoryAddr.Org))+"/"+url.PathEscape(repositoryAddr.Name), &repoResponse); err != nil {
		var statusCodeError *InvalidStatusCodeError
		if errors.As(err, &statusCodeError) && statusCodeError.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type InvalidStatusCodeError struct {
	StatusCode int
}

func (i InvalidStatusCodeError) Error() string {
	return "Invalid status code: " + strconv.Itoa(i.StatusCode)
}
