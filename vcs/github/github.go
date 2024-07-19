package github

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/opentofu/libregistry/vcs"
)

// New creates a new GitHub VCS client.
func New(
	token string,
	httpClient *http.Client,
) (vcs.Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
		transport := http.DefaultTransport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
		httpClient.Transport = transport
	}

	return &github{
		token:      token,
		httpClient: httpClient,
	}, nil
}

type github struct {
	token      string
	httpClient *http.Client
}

func (g github) ParseRepositoryAddr(ref string) (vcs.RepositoryAddr, error) {
	ref = strings.TrimPrefix(ref, "github.com/")
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) != 2 {
		return vcs.RepositoryAddr{}, &vcs.InvalidRepositoryAddrError{
			RepositoryAddr: ref,
		}
	}
	result := vcs.RepositoryAddr{
		OrganizationAddr: vcs.OrganizationAddr{Org: parts[0]},
		Name:             parts[1],
	}
	return result, result.Validate()
}

type rss struct {
	Entry []struct {
		ID    string `title:"id"`
		Title string `xml:"title"`
	} `xml:"entry"`
}

func (g github) ListLatestVersions(ctx context.Context, repository vcs.RepositoryAddr) ([]string, error) {
	if err := repository.Validate(); err != nil {
		return nil, err
	}
	rssURL := "https://github.com/" + url.PathEscape(repository.Org) + "/" + url.PathEscape(repository.Name) + "/tags.atom"
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

	result := make([]string, len(response.Entry))
	for i, entry := range response.Entry {
		result[i] = entry.Title
	}
	// TODO handle incorrectly named version.
	return result, nil
}

func (g github) ListAllVersions(ctx context.Context, repository vcs.RepositoryAddr) ([]string, error) {
	type tag struct {
		Name string `json:"name"`
	}

	if err := repository.Validate(); err != nil {
		return nil, err
	}
	rssURL := "https://api.github.com/repos/" + url.PathEscape(repository.Org) + "/" + url.PathEscape(repository.Name) + "/releases"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rssURL, nil)
	if err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: fmt.Errorf("invalid HTTP request (%w)", err),
		}
	}
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
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

	decoder := json.NewDecoder(resp.Body)

	var response []tag
	if err := decoder.Decode(&response); err != nil {
		return nil, &vcs.RequestFailedError{
			Cause: fmt.Errorf("failed to decode response (%w)", err),
		}
	}

	result := make([]string, len(response))
	for i, item := range response {
		result[i] = item.Name
	}
	return result, nil
}

func (g github) ListAssets(ctx context.Context, repository vcs.RepositoryAddr, version string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (g github) DownloadAsset(ctx context.Context, repository vcs.RepositoryAddr, version string, asset string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g github) HasPermission(ctx context.Context, organization vcs.OrganizationAddr) (bool, error) {
	//TODO implement me
	panic("implement me")
}
