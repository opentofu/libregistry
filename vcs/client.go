// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package vcs

import (
	"context"
	"io/fs"
)

// Client describes a VCS client.
type Client interface {
	// ParseRepositoryAddr parses the repository address from a string.
	ParseRepositoryAddr(ref string) (RepositoryAddr, error)

	// GetRepositoryInfo returns the information from the VCS API.
	GetRepositoryInfo(ctx context.Context, repository RepositoryAddr) (RepositoryInfo, error)

	// ListLatestTags returns the last few tags in the VCS system. This is a lightweight call and
	// may need to be supplemented by a call to ListAllTags.
	//
	// Caution! This function MAY perform a checkout, which may place an exclusive lock on the checkout directory.
	// Do not call this function while you have a working copy checked out!
	ListLatestTags(ctx context.Context, repository RepositoryAddr) ([]Version, error)

	// ListAllTags returns a list of all tags in the repository. Whenever possible, prefer
	// ListLatestTags instead since this call may be heavily rate limited.
	//
	// Caution! This function MAY perform a checkout, which may place an exclusive lock on the checkout directory.
	// Do not call this function while you have a working copy checked out!
	ListAllTags(ctx context.Context, repository RepositoryAddr) ([]Version, error)

	// GetTagVersion returns the full version information for a tag.
	//
	// Caution! This function MAY perform a checkout, which may place an exclusive lock on the checkout directory.
	// Do not call this function while you have a working copy checked out!
	GetTagVersion(ctx context.Context, repository RepositoryAddr, version VersionNumber) (Version, error)

	// ListLatestReleases returns the last few releases in the VCS system. This is a lightweight call and
	// may need to be supplemented by a call to ListAllReleases.
	//
	// Caution! This function MAY perform a checkout, which may place an exclusive lock on the checkout directory.
	// Do not call this function while you have a working copy checked out!
	ListLatestReleases(ctx context.Context, repository RepositoryAddr) ([]Version, error)

	// ListAllReleases returns a list of all releases in the repository. Whenever possible, prefer
	// ListLatestReleases instead since this call may be heavily rate limited.
	//
	// Caution! This function MAY perform a checkout, which may place an exclusive lock on the checkout directory.
	// Do not call this function while you have a working copy checked out!
	ListAllReleases(ctx context.Context, repository RepositoryAddr) ([]Version, error)

	// ListAssets lists all binary assets for a release of a repository.
	ListAssets(ctx context.Context, repository RepositoryAddr, version VersionNumber) ([]AssetName, error)

	// DownloadAsset downloads a given asset from a release in a repository.
	DownloadAsset(ctx context.Context, repository RepositoryAddr, version VersionNumber, asset AssetName) ([]byte, error)

	// HasPermission returns true if the user has permission to act on behalf of an organization.
	HasPermission(ctx context.Context, username Username, organization OrganizationAddr) (bool, error)

	// Checkout clones/checks out a working copy of the given repository at a given version for accessing the files
	// in the repository. The caller is responsible for calling Close on the WorkingCopy when finished to allow a
	// cleanup or release any locks.
	//
	// Note that the implementation may limit the concurrent use of a repository to a single WorkingCopy at a time in
	// order to adhere to any rate limits the VCS system may impose.
	Checkout(ctx context.Context, repository RepositoryAddr, version VersionNumber) (WorkingCopy, error)

	// GetRepositoryBrowseURL returns the web address the repository can be viewed at. The implementation may return
	// a *NoWebAccessError if the VCS system does not support accessing files via the web.
	GetRepositoryBrowseURL(ctx context.Context, repository RepositoryAddr) (string, error)

	// GetVersionBrowseURL returns the web address a specific version can be viewed at. The implementation may return
	// a *NoWebAccessError if the VCS system does not support accessing files via the web.
	GetVersionBrowseURL(ctx context.Context, repository RepositoryAddr, version VersionNumber) (string, error)

	// GetFileViewURL determines the URL a file in a specific version can be viewed at. The existence of the file in
	// the version is not verified. The implementation may return a *NoWebAccessError if the VCS system does not support
	// accessing files via the web.
	GetFileViewURL(ctx context.Context, repository RepositoryAddr, version VersionNumber, file string) (string, error)
}

type WorkingCopy interface {
	fs.ReadDirFS
	// RawDirectory returns the underlying raw directory of the working copy. This should be used with care as any open
	// file descriptors may cause the Close() call to fail. Modifications to this directory should also be avoided.
	// This call may return an error if raw directory access is not supported.
	RawDirectory() (string, error)
	// Client returns the VCS client this working copy belongs to.
	Client() Client
	// Repository returns the repository for the current working copy.
	Repository() RepositoryAddr
	// Version returns the version number the working copy is checked out at.
	Version() VersionNumber
	// Close cleans up the working copy and releases the lock on it.
	Close() error
}

type RepositoryInfo struct {
	Description string `json:"description"`
	// Popularity indicates how popular (stars, etc.) the repository is.
	Popularity int `json:"popularity"`
	// ForkOf indicates that this repository is a fork or copy of another repository. This is empty if the
	// repository is not a known fork/copy.
	ForkOf *RepositoryAddr `json:"fork_of,omitempty"`
	// ForkCount exposes the amount of copies/forks present in the VCS.
	ForkCount int
}
