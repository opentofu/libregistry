// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package fakevcs

import (
	"context"
	"strings"

	"github.com/opentofu/libregistry/vcs"
)

type inMemoryVCS struct {
	users         map[string]struct{}
	organizations map[vcs.OrganizationAddr]*org
}

func (i *inMemoryVCS) ParseRepositoryAddr(ref string) (vcs.RepositoryAddr, error) {
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) != 2 {
		return vcs.RepositoryAddr{}, &vcs.InvalidRepositoryAddrError{
			RepositoryAddr: ref,
			Cause:          nil,
		}
	}
	addr := vcs.RepositoryAddr{
		OrganizationAddr: vcs.OrganizationAddr{
			Org: parts[0],
		},
		Name: parts[1],
	}
	return addr, addr.Validate()
}

func (i *inMemoryVCS) ListLatestVersions(ctx context.Context, repositoryAddr vcs.RepositoryAddr) ([]string, error) {
	versions, err := i.ListAllVersions(ctx, repositoryAddr)
	if len(versions) > 5 {
		versions = versions[:5]
	}
	return versions, err
}

func (i *inMemoryVCS) ListAllVersions(_ context.Context, repositoryAddr vcs.RepositoryAddr) ([]string, error) {
	org, ok := i.organizations[repositoryAddr.OrganizationAddr]
	if !ok {
		return nil, &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	repo, ok := org.repositories[repositoryAddr]
	if !ok {
		return nil, &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}

	result := make([]string, len(repo.versions))
	for i, ver := range repo.versions {
		result[i] = ver.name
	}
	return result, nil
}

func (i *inMemoryVCS) ListAssets(_ context.Context, repositoryAddr vcs.RepositoryAddr, version string) ([]string, error) {
	org, ok := i.organizations[repositoryAddr.OrganizationAddr]
	if !ok {
		return nil, &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	repo, ok := org.repositories[repositoryAddr]
	if !ok {
		return nil, &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	for _, ver := range repo.versions {
		if ver.name == version {
			result := make([]string, len(ver.assets))
			i := 0
			for name := range ver.assets {
				result[i] = name
				i++
			}
			return result, nil
		}
	}
	return nil, &vcs.VersionNotFoundError{
		RepositoryAddr: repositoryAddr,
		Version:        version,
	}
}

func (i *inMemoryVCS) DownloadAsset(_ context.Context, repositoryAddr vcs.RepositoryAddr, version string, asset string) ([]byte, error) {
	org, ok := i.organizations[repositoryAddr.OrganizationAddr]
	if !ok {
		return nil, &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	repo, ok := org.repositories[repositoryAddr]
	if !ok {
		return nil, &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	for _, ver := range repo.versions {
		if ver.name == version {
			if assetData, ok := ver.assets[asset]; ok {
				return assetData, nil
			}
		}
	}
	return nil, &vcs.AssetNotFoundError{
		RepositoryAddr: repositoryAddr,
		Version:        version,
		Asset:          asset,
	}
}

func (i *inMemoryVCS) HasPermission(_ context.Context, username string, organization vcs.OrganizationAddr) (bool, error) {
	org, ok := i.organizations[organization]
	if !ok {
		return false, &vcs.OrganizationNotFoundError{
			OrganizationAddr: organization,
		}
	}
	_, ok = org.users[username]
	return ok, nil
}

func (i *inMemoryVCS) CreateOrganization(organization vcs.OrganizationAddr) error {
	if _, ok := i.organizations[organization]; ok {
		return &OrganizationAlreadyExistsError{organization}
	}
	i.organizations[organization] = &org{
		users:        map[string]struct{}{},
		repositories: map[vcs.RepositoryAddr]*repository{},
	}
	return nil
}

func (i *inMemoryVCS) CreateRepository(repositoryAddr vcs.RepositoryAddr) error {
	if _, ok := i.organizations[repositoryAddr.OrganizationAddr]; !ok {
		return &vcs.OrganizationNotFoundError{
			OrganizationAddr: repositoryAddr.OrganizationAddr,
		}
	}
	if _, ok := i.organizations[repositoryAddr.OrganizationAddr].repositories[repositoryAddr]; ok {
		return &RepositoryAlreadyExistsError{
			repositoryAddr,
		}
	}
	i.organizations[repositoryAddr.OrganizationAddr].repositories[repositoryAddr] = &repository{}
	return nil
}

func (i *inMemoryVCS) CreateVersion(repositoryAddr vcs.RepositoryAddr, versionName string) error {
	org, ok := i.organizations[repositoryAddr.OrganizationAddr]
	if !ok {
		return &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	repo, ok := org.repositories[repositoryAddr]
	if !ok {
		return &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	for _, ver := range repo.versions {
		if ver.name == versionName {
			return &VersionAlreadyExistsError{
				repositoryAddr,
				versionName,
			}
		}
	}
	repo.versions = append([]version{
		{
			name:   versionName,
			assets: map[string][]byte{},
		},
	}, repo.versions...)
	return nil
}

func (i *inMemoryVCS) AddAsset(repositoryAddr vcs.RepositoryAddr, versionName string, assetName string, assetData []byte) error {
	org, ok := i.organizations[repositoryAddr.OrganizationAddr]
	if !ok {
		return &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	repo, ok := org.repositories[repositoryAddr]
	if !ok {
		return &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	for _, ver := range repo.versions {
		if ver.name == versionName {
			if _, ok := ver.assets[assetName]; ok {
				return AssetAlreadyExistsError{
					repositoryAddr,
					versionName,
					assetName,
				}
			}
			ver.assets[assetName] = assetData
			return nil
		}
	}

	return &vcs.VersionNotFoundError{
		RepositoryAddr: repositoryAddr,
		Version:        versionName,
		Cause:          nil,
	}
}

func (i *inMemoryVCS) AddUser(username string) error {
	if _, ok := i.users[username]; ok {
		return &UserAlreadyExistsError{
			username,
		}
	}
	i.users[username] = struct{}{}
	return nil
}

func (i *inMemoryVCS) AddMember(organizationAddr vcs.OrganizationAddr, username string) error {
	if _, ok := i.users[username]; !ok {
		return &UserNotFoundError{
			username,
		}
	}
	org, ok := i.organizations[organizationAddr]
	if !ok {
		return &vcs.OrganizationNotFoundError{
			OrganizationAddr: organizationAddr,
		}
	}
	org.users[username] = struct{}{}
	return nil
}
