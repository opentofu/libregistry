// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package libregistry

import (
	"context"
	"errors"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/module"
	"github.com/opentofu/libregistry/vcs"
)

func (m api) UpdateModule(ctx context.Context, moduleAddr module.Addr) error {
	if err := moduleAddr.Validate(); err != nil {
		return &ModuleUpdateFailedError{
			moduleAddr,
			err,
		}
	}

	moduleMetadata, err := m.dataAPI.GetModule(ctx, moduleAddr)
	if err != nil {
		var notFoundError *metadata.ModuleNotFoundError
		if !errors.As(err, &notFoundError) {
			return &ModuleUpdateFailedError{
				moduleAddr,
				err,
			}
		}
		moduleMetadata = module.Metadata{}
	}

	previousSize := len(moduleMetadata.Versions)
	tags, err := m.vcsClient.ListLatestTags(ctx, getModuleRepo(moduleAddr))
	if err != nil {
		return &ModuleUpdateFailedError{
			moduleAddr,
			err,
		}
	}
	var newVersions module.VersionList
	for _, tag := range tags {
		ver, err := module.VersionFromVCS(tag)
		if err != nil {
			continue
		}
		newVersions = append(newVersions, module.Version{
			Version: ver,
		})
	}
	moduleMetadata.Versions = moduleMetadata.Versions.Merge(newVersions)

	if len(moduleMetadata.Versions) == previousSize+len(tags) {
		// No overlap found, do the full query:
		tags, err = m.vcsClient.ListAllTags(ctx, getModuleRepo(moduleAddr))
		if err != nil {
			return &ModuleUpdateFailedError{
				moduleAddr,
				err,
			}
		}
		newVersions = nil
		for _, tag := range tags {
			ver, err := module.VersionFromVCS(tag)
			if err != nil {
				continue
			}
			newVersions = append(newVersions, module.Version{
				Version: ver,
			})
		}
		moduleMetadata.Versions = newVersions
	}

	if err := m.dataAPI.PutModule(ctx, moduleAddr, moduleMetadata); err != nil {
		return &ModuleAddFailedError{
			moduleAddr,
			err,
		}
	}
	return nil
}

func getModuleRepo(module module.Addr) vcs.RepositoryAddr {
	return vcs.RepositoryAddr{
		Org:  vcs.OrganizationAddr(module.Namespace),
		Name: "terraform-" + module.TargetSystem + "-" + module.Name,
	}
}
