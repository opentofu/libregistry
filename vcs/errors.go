// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package vcs

type InvalidRepositoryAddrError struct {
	RepositoryAddr string
	Cause          error
}

func (r InvalidRepositoryAddrError) Error() string {
	if r.Cause != nil {
		return "Failed to parse repository address: " + r.RepositoryAddr + " (" + r.Cause.Error() + ")"
	}
	return "Failed to parse repository address: " + r.RepositoryAddr
}

func (r InvalidRepositoryAddrError) Unwrap() error {
	return r.Cause
}

type RequestFailedError struct {
	Cause error
}

func (r RequestFailedError) Error() string {
	return "VCS request failed: " + r.Cause.Error()
}

func (r RequestFailedError) Unwrap() error {
	return r.Cause
}

type OrganizationNotFoundError struct {
	OrganizationAddr OrganizationAddr
	Cause            error
}

func (r OrganizationNotFoundError) Error() string {
	if r.Cause != nil {
		return "Organization not found: " + r.OrganizationAddr.String() + " (" + r.Cause.Error() + ")"
	}
	return "Organization not found: " + r.OrganizationAddr.String()
}

func (r OrganizationNotFoundError) Unwrap() error {
	return r.Cause
}

type RepositoryNotFoundError struct {
	RepositoryAddr RepositoryAddr
	Cause          error
}

func (r RepositoryNotFoundError) Error() string {
	if r.Cause != nil {
		return "Repository not found: " + r.RepositoryAddr.String() + " (" + r.Cause.Error() + ")"
	}
	return "Repository not found: " + r.RepositoryAddr.String()
}

func (r RepositoryNotFoundError) Unwrap() error {
	return r.Cause
}

type VersionNotFoundError struct {
	RepositoryAddr RepositoryAddr
	Version        string
	Cause          error
}

func (v VersionNotFoundError) Error() string {
	if v.Cause != nil {
		return "Version " + v.Version + " not found: in repository" + v.RepositoryAddr.String() + " (" + v.Cause.Error() + ")"
	}
	return "Version " + v.Version + " not found: in repository" + v.RepositoryAddr.String()
}

func (v VersionNotFoundError) Unwrap() error {
	return v.Cause
}

type AssetNotFoundError struct {
	RepositoryAddr RepositoryAddr
	Version        string
	Asset          string
	Cause          error
}

func (a AssetNotFoundError) Error() string {
	if a.Cause != nil {
		return "Asset " + a.Asset + " not found in version " + a.Version + " of repository" + a.RepositoryAddr.String() + " (" + a.Cause.Error() + ")"
	}
	return "Asset " + a.Asset + " not found in version " + a.Version + " of repository" + a.RepositoryAddr.String()
}

func (a AssetNotFoundError) Unwrap() error {
	return a.Cause
}
