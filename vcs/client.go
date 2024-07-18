package vcs

import (
	"context"
)

// Client describes a VCS client.
type Client interface {
	// ParseRepositoryAddr parses the repository address from a string.
	ParseRepositoryAddr(ref string) (RepositoryAddr, error)

	// ListVersions returns all versions (e.g. tags) in the VCS system.
	ListVersions(ctx context.Context, repository RepositoryAddr) ([]string, error)

	// ListAssets lists all binary assets for a version of a repository.
	ListAssets(ctx context.Context, repository RepositoryAddr, version string) ([]string, error)

	// DownloadAsset downloads a given asset from a release in a repository.
	DownloadAsset(ctx context.Context, repository RepositoryAddr, version string, asset string) ([]byte, error)

	// HasPermission returns true if the user has permission to act on behalf of an organization.
	HasPermission(ctx context.Context, organization OrganizationAddr) (bool, error)
}
