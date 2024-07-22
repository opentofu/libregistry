// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package vcs

import (
	"context"
)

// Client describes a VCS client.
type Client interface {
	// ParseRepositoryAddr parses the repository address from a string.
	ParseRepositoryAddr(ref string) (RepositoryAddr, error)

	// ListLatestTags returns the last few tags in the VCS system. This is a lightweight call and
	// may need to be supplemented by a call to ListAllTags.
	ListLatestTags(ctx context.Context, repository RepositoryAddr) ([]Version, error)

	// ListAllTags returns a list of all tags in the repository. Whenever possible, prefer
	// ListLatestTags instead since this call may be heavily rate limited.
	ListAllTags(ctx context.Context, repository RepositoryAddr) ([]Version, error)

	// ListLatestReleases returns the last few releases in the VCS system. This is a lightweight call and
	// may need to be supplemented by a call to ListAllReleases.
	ListLatestReleases(ctx context.Context, repository RepositoryAddr) ([]Version, error)

	// ListAllReleases returns a list of all releases in the repository. Whenever possible, prefer
	// ListLatestReleases instead since this call may be heavily rate limited.
	ListAllReleases(ctx context.Context, repository RepositoryAddr) ([]Version, error)

	// ListAssets lists all binary assets for a release of a repository.
	ListAssets(ctx context.Context, repository RepositoryAddr, version Version) ([]AssetName, error)

	// DownloadAsset downloads a given asset from a release in a repository.
	DownloadAsset(ctx context.Context, repository RepositoryAddr, version Version, asset AssetName) ([]byte, error)

	// HasPermission returns true if the user has permission to act on behalf of an organization.
	HasPermission(ctx context.Context, username Username, organization OrganizationAddr) (bool, error)
}
