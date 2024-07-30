// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package fakevcs

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/opentofu/libregistry/vcs"
)

type inMemoryVCS struct {
	config        Config
	users         map[vcs.Username]struct{}
	organizations map[vcs.OrganizationAddr]*org
}

func (i *inMemoryVCS) GetRepositoryInfo(_ context.Context, repositoryAddr vcs.RepositoryAddr) (vcs.RepositoryInfo, error) {
	if err := repositoryAddr.Validate(); err != nil {
		return vcs.RepositoryInfo{}, err
	}

	org, ok := i.organizations[repositoryAddr.Org]
	if !ok {
		return vcs.RepositoryInfo{}, &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}
	repo, ok := org.repositories[repositoryAddr]
	if !ok {
		return vcs.RepositoryInfo{}, &vcs.RepositoryNotFoundError{
			RepositoryAddr: repositoryAddr,
		}
	}

	return repo.info, nil
}

func (i *inMemoryVCS) ListLatestReleases(ctx context.Context, repository vcs.RepositoryAddr) ([]vcs.Version, error) {
	return i.ListAllReleases(ctx, repository)
}

func (i *inMemoryVCS) ListAllReleases(_ context.Context, repositoryAddr vcs.RepositoryAddr) ([]vcs.Version, error) {
	if err := repositoryAddr.Validate(); err != nil {
		return nil, err
	}

	org, ok := i.organizations[repositoryAddr.Org]
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

	result := make([]vcs.Version, len(repo.versions))
	for i, ver := range repo.versions {
		result[i] = vcs.Version{
			VersionNumber: ver.name,
			Created:       ver.created,
		}
	}
	return result, nil
}

func (i *inMemoryVCS) ParseRepositoryAddr(ref string) (vcs.RepositoryAddr, error) {
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) != 2 {
		return vcs.RepositoryAddr{}, &vcs.InvalidRepositoryAddrError{
			RepositoryString: ref,
		}
	}
	addr := vcs.RepositoryAddr{
		Org:  vcs.OrganizationAddr(parts[0]),
		Name: parts[1],
	}
	return addr, addr.Validate()
}

func (i *inMemoryVCS) ListLatestTags(ctx context.Context, repositoryAddr vcs.RepositoryAddr) ([]vcs.Version, error) {
	versions, err := i.ListAllTags(ctx, repositoryAddr)
	if len(versions) > 5 {
		versions = versions[:5]
	}
	return versions, err
}

func (i *inMemoryVCS) ListAllTags(_ context.Context, repositoryAddr vcs.RepositoryAddr) ([]vcs.Version, error) {
	if err := repositoryAddr.Validate(); err != nil {
		return nil, err
	}

	org, ok := i.organizations[repositoryAddr.Org]
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

	result := make([]vcs.Version, len(repo.versions))
	for i, ver := range repo.versions {
		result[i] = vcs.Version{
			VersionNumber: ver.name,
			Created:       ver.created,
		}
	}
	return result, nil
}

func (i *inMemoryVCS) ListAssets(_ context.Context, repositoryAddr vcs.RepositoryAddr, version vcs.VersionNumber) ([]vcs.AssetName, error) {
	if err := repositoryAddr.Validate(); err != nil {
		return nil, err
	}
	if err := version.Validate(); err != nil {
		return nil, err
	}
	org, ok := i.organizations[repositoryAddr.Org]
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
			result := make([]vcs.AssetName, len(ver.assets))
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

func (i *inMemoryVCS) DownloadAsset(_ context.Context, repositoryAddr vcs.RepositoryAddr, version vcs.VersionNumber, asset vcs.AssetName) ([]byte, error) {
	if err := repositoryAddr.Validate(); err != nil {
		return nil, err
	}
	if err := version.Validate(); err != nil {
		return nil, err
	}
	if err := asset.Validate(); err != nil {
		return nil, err
	}
	org, ok := i.organizations[repositoryAddr.Org]
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

func (i *inMemoryVCS) HasPermission(_ context.Context, username vcs.Username, organization vcs.OrganizationAddr) (bool, error) {
	if err := organization.Validate(); err != nil {
		return false, err
	}
	if err := username.Validate(); err != nil {
		return false, err
	}
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
	if err := organization.Validate(); err != nil {
		return err
	}
	if _, ok := i.organizations[organization]; ok {
		return &OrganizationAlreadyExistsError{organization}
	}
	i.organizations[organization] = &org{
		users:        map[vcs.Username]struct{}{},
		repositories: map[vcs.RepositoryAddr]*repository{},
	}
	return nil
}

func (i *inMemoryVCS) CreateRepository(repositoryAddr vcs.RepositoryAddr, repositoryInfo vcs.RepositoryInfo) error {
	if err := repositoryAddr.Validate(); err != nil {
		return err
	}
	if _, ok := i.organizations[repositoryAddr.Org]; !ok {
		return &vcs.OrganizationNotFoundError{
			OrganizationAddr: repositoryAddr.Org,
		}
	}
	if _, ok := i.organizations[repositoryAddr.Org].repositories[repositoryAddr]; ok {
		return &RepositoryAlreadyExistsError{
			repositoryAddr,
		}
	}
	i.organizations[repositoryAddr.Org].repositories[repositoryAddr] = &repository{}
	return nil
}

func (i *inMemoryVCS) CreateVersion(repositoryAddr vcs.RepositoryAddr, versionName vcs.VersionNumber, contents fs.ReadDirFS) error {
	if err := repositoryAddr.Validate(); err != nil {
		return err
	}
	if err := versionName.Validate(); err != nil {
		return err
	}
	org, ok := i.organizations[repositoryAddr.Org]
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
			name:     versionName,
			created:  i.config.TimeSource(),
			assets:   map[vcs.AssetName][]byte{},
			contents: contents,
		},
	}, repo.versions...)
	return nil
}

func (i *inMemoryVCS) AddAsset(repositoryAddr vcs.RepositoryAddr, versionName vcs.VersionNumber, assetName vcs.AssetName, assetData []byte) error {
	if err := repositoryAddr.Validate(); err != nil {
		return err
	}
	if err := versionName.Validate(); err != nil {
		return err
	}
	if err := assetName.Validate(); err != nil {
		return err
	}
	org, ok := i.organizations[repositoryAddr.Org]
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

func (i *inMemoryVCS) AddUser(username vcs.Username) error {
	if _, ok := i.users[username]; ok {
		return &UserAlreadyExistsError{
			username,
		}
	}
	i.users[username] = struct{}{}
	return nil
}

func (i *inMemoryVCS) AddMember(organizationAddr vcs.OrganizationAddr, username vcs.Username) error {
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

func (i *inMemoryVCS) Checkout(ctx context.Context, repositoryAddr vcs.RepositoryAddr, version vcs.VersionNumber) (vcs.WorkingCopy, error) {
	if err := repositoryAddr.Validate(); err != nil {
		return nil, err
	}
	if err := version.Validate(); err != nil {
		return nil, err
	}
	org, ok := i.organizations[repositoryAddr.Org]
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
			return &workingCopy{
				ver.contents,
			}, nil
		}
	}
	return nil, &vcs.VersionNotFoundError{
		RepositoryAddr: repositoryAddr,
		Version:        version,
		Cause:          nil,
	}
}

type workingCopy struct {
	fs.ReadDirFS
}

func (w workingCopy) RawDirectory() (string, error) {
	return "", fmt.Errorf("raw directory access is not supported for the fake VCS")
}

func (w workingCopy) Close() error {
	return nil
}
