package vcs

import (
	"context"
)

// Client describes a VCS client.
type Client interface {
	// ParseRepositoryAddr parses the repository address from a string.
	ParseRepositoryAddr(ref string) (RepositoryAddr, error)

	// ListLatestVersions returns the last few versions (e.g. tags) in the VCS system. This is a lightweight call and
	// may need to be supplemented by a call to ListAllVersions.
	ListLatestVersions(ctx context.Context, repository RepositoryAddr) ([]string, error)

	// ListAllVersions returns a list of all versions (e.g. tags) in the repository. Whenever possible, prefer
	// ListLatestVersions instead since this call may be heavily rate limited.
	ListAllVersions(ctx context.Context, repository RepositoryAddr) ([]string, error)

	// ListAssets lists all binary assets for a version of a repository.
	ListAssets(ctx context.Context, repository RepositoryAddr, version string) ([]string, error)

	// DownloadAsset downloads a given asset from a release in a repository.
	DownloadAsset(ctx context.Context, repository RepositoryAddr, version string, asset string) ([]byte, error)

	// HasPermission returns true if the user has permission to act on behalf of an organization.
	HasPermission(ctx context.Context, username string, organization OrganizationAddr) (bool, error)
}
