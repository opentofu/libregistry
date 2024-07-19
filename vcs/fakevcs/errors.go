// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package fakevcs

import (
	"github.com/opentofu/libregistry/vcs"
)

type OrganizationAlreadyExistsError struct {
	OrganizationAddr vcs.OrganizationAddr
}

func (r OrganizationAlreadyExistsError) Error() string {
	return "Organization already exists: " + r.OrganizationAddr.String()
}

type RepositoryAlreadyExistsError struct {
	RepositoryAddr vcs.RepositoryAddr
}

func (r RepositoryAlreadyExistsError) Error() string {
	return "Repository already exists: " + r.RepositoryAddr.String()
}

type VersionAlreadyExistsError struct {
	RepositoryAddr vcs.RepositoryAddr
	Version        string
}

func (v VersionAlreadyExistsError) Error() string {
	return "Version " + v.Version + " already exists in repository " + v.RepositoryAddr.String()
}

type AssetAlreadyExistsError struct {
	RepositoryAddr vcs.RepositoryAddr
	Version        string
	Asset          string
}

func (a AssetAlreadyExistsError) Error() string {
	return "Asset + " + a.Asset + " already exists for version " + a.Version + " on repository " + a.RepositoryAddr.String()
}

type UserAlreadyExistsError struct {
	Username string
}

func (u UserAlreadyExistsError) Error() string {
	return "User " + u.Username + " already exists."
}

type UserNotFoundError struct {
	Username string
}

func (u UserNotFoundError) Error() string {
	return "User " + u.Username + " not found."
}
