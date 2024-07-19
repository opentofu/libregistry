package fakevcs

import (
	"github.com/opentofu/libregistry/vcs"
)

// New creates a fake, in-memory VCSClient implementation for testing use.
func New() VCSClient {
	return &inMemoryVCS{
		users:         map[string]struct{}{},
		organizations: map[vcs.OrganizationAddr]*org{},
	}
}

type VCSClient interface {
	vcs.Client

	CreateOrganization(organization vcs.OrganizationAddr) error
	CreateRepository(repository vcs.RepositoryAddr) error
	CreateVersion(repository vcs.RepositoryAddr, version string) error
	AddAsset(repository vcs.RepositoryAddr, version string, name string, data []byte) error
	AddUser(username string) error
	AddMember(organization vcs.OrganizationAddr, username string) error
}
