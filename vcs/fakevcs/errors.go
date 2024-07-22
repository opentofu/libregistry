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
	Version        vcs.Version
}

func (v VersionAlreadyExistsError) Error() string {
	return "Version " + string(v.Version) + " already exists in repository " + v.RepositoryAddr.String()
}

type AssetAlreadyExistsError struct {
	RepositoryAddr vcs.RepositoryAddr
	Version        vcs.Version
	Asset          vcs.AssetName
}

func (a AssetAlreadyExistsError) Error() string {
	return "Asset + " + string(a.Asset) + " already exists for version " + string(a.Version) + " on repository " + a.RepositoryAddr.String()
}

type UserAlreadyExistsError struct {
	Username vcs.Username
}

func (u UserAlreadyExistsError) Error() string {
	return "User " + string(u.Username) + " already exists."
}

type UserNotFoundError struct {
	Username vcs.Username
}

func (u UserNotFoundError) Error() string {
	return "User " + string(u.Username) + " not found."
}
