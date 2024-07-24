// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package fakevcs

import (
	"github.com/opentofu/libregistry/vcs"
)

// New creates a fake, in-memory VCSClient implementation for testing use.
func New() VCSClient {
	return &inMemoryVCS{
		users:         map[vcs.Username]struct{}{},
		organizations: map[vcs.OrganizationAddr]*org{},
	}
}

type VCSClient interface {
	vcs.Client

	CreateOrganization(organization vcs.OrganizationAddr) error
	CreateRepository(repository vcs.RepositoryAddr) error
	CreateVersion(repository vcs.RepositoryAddr, version vcs.Version) error
	AddAsset(repository vcs.RepositoryAddr, version vcs.Version, name vcs.AssetName, data []byte) error
	AddUser(username vcs.Username) error
	AddMember(organization vcs.OrganizationAddr, username vcs.Username) error
}
