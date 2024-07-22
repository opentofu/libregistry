// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/opentofu/libregistry/logger"
	"github.com/opentofu/libregistry/vcs"
)

// New creates a new GitHub VCS client.
func New(
	token string,
	httpClient *http.Client,
	log logger.Logger,
) (vcs.Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
		transport := http.DefaultTransport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
		httpClient.Transport = transport
	}

	if log == nil {
		log = logger.NewNoopLogger()
	}

	return &github{
		token:      token,
		httpClient: httpClient,
		logger:     log,
	}, nil
}

type github struct {
	token      string
	httpClient *http.Client
	logger     logger.Logger
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
		ID    string `title:"id"`
		Title string `xml:"title"`
	} `xml:"entry"`
}

func (g github) ListLatestTags(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	logger.LogTrace(ctx, g.logger, "Requesting latest tags for repository %s...", repository)
	return g.listLatest(ctx, repository, "tags.atom")
}

func (g github) ListLatestReleases(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	logger.LogTrace(ctx, g.logger, "Requesting latest releases for repository %s...", repository)
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
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: err,
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		return nil, &vcs.RequestFailedError{
			Cause: fmt.Errorf("invalid status code: %d", resp.StatusCode),
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
		version := vcs.Version(entry.Title)
		err = version.Validate()
		if err != nil {
			g.logger.Debug(ctx, "Skipping invalid version %s when querying %s in repository %s", version, file, repository)
		} else {
			result = append(result, version)
		}
	}
	return result, nil
}

func (g github) ListAllTags(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	logger.LogTrace(ctx, g.logger, "Requesting all tags for repository %s...", repository)
	return g.listAll(ctx, repository, "tag")
}

func (g github) ListAllReleases(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	logger.LogTrace(ctx, g.logger, "Requesting all releases for repository %s...", repository)
	return g.listAll(ctx, repository, "release")
}

func (g github) listAll(ctx context.Context, repository vcs.RepositoryAddr, itemType string) ([]vcs.Version, error) {
	type responseItem struct {
		Name vcs.Version `json:"name"`
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
			g.logger.Debug(ctx, "Skipping invalid %s %s in repository %s", itemType, item.Name, repository)
		} else {
			result = append(result, item.Name)
		}
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
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}
	logger.LogTrace(ctx, g.logger, "Sending GET request to %s...", url)
	resp, err := g.httpClient.Do(req)
	if err != nil {
		logger.LogTrace(ctx, g.logger, "GET request to %s failed (%v)", url, err)
		return &vcs.RequestFailedError{
			Cause: err,
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	logger.LogTrace(ctx, g.logger, "GET request to %s returned status code %d", url, resp.StatusCode)
	if resp.StatusCode != 200 {
		return &vcs.RequestFailedError{
			Cause: &InvalidStatusCodeError{resp.StatusCode},
		}
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&response); err != nil {
		g.logger.Warn(ctx, "GitHub returned an invalid JSON when requesting %s (%v)", url, err)
		return &vcs.RequestFailedError{
			Cause: fmt.Errorf("failed to decode response (%w)", err),
		}
	}

	return nil
}

func (g github) ListAssets(ctx context.Context, repository vcs.RepositoryAddr, version vcs.Version) ([]vcs.AssetName, error) {
	logger.LogTrace(ctx, g.logger, "Listing assets for repository %s version %s", repository, version)

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
			g.logger.Debug(ctx, "Skipping invalid asset named %s in repository %s release %s", asset.Name, repository, version)
		} else {
			result = append(result, asset.Name)
		}
	}
	return result, nil
}

func (g github) DownloadAsset(ctx context.Context, repository vcs.RepositoryAddr, version vcs.Version, asset vcs.AssetName) ([]byte, error) {
	if err := repository.Validate(); err != nil {
		return nil, err
	}
	if err := version.Validate(); err != nil {
		return nil, err
	}
	if err := asset.Validate(); err != nil {
		return nil, err
	}
	logger.LogTrace(ctx, g.logger, "Listing asset %s for repository %s version %s", asset, repository, version)
	assetURL := "https://api.github.com/repos/" + url.PathEscape(string(repository.Org)) + "/" + url.PathEscape(repository.Name) + "/releases/download/" + url.PathEscape(string(version)) + "/" + url.PathEscape(string(asset))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, assetURL, nil)
	if err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: fmt.Errorf("invalid HTTP request (%w)", err),
		}
	}
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}
	logger.LogTrace(ctx, g.logger, "Sending GET request to %s...", assetURL)
	resp, err := g.httpClient.Do(req)
	if err != nil {
		logger.LogTrace(ctx, g.logger, "GET request to %s failed (%v)", assetURL, err)
		return nil, &vcs.RequestFailedError{
			Cause: err,
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	logger.LogTrace(ctx, g.logger, "GET request to %s returned status code %d", assetURL, resp.StatusCode)
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
	logger.LogTrace(ctx, g.logger, "Checking if user %s has permissions for the organization %s...", username, organization)
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
		if strings.ToLower(member.Login) == strings.ToLower(string(username)) {
			return true, nil
		}
	}
	return false, nil
}

type InvalidStatusCodeError struct {
	StatusCode int
}

func (i InvalidStatusCodeError) Error() string {
	return "Invalid status code: " + strconv.Itoa(i.StatusCode)
}
