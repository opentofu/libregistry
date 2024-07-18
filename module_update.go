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

	tags, err := m.vcsClient.ListVersions(ctx, getModuleRepo(moduleAddr))
	if err != nil {
		return &ModuleAddFailedError{
			moduleAddr,
			err,
		}
	}
	var newVersions module.VersionList
	for _, tag := range tags {
		newVersions = append(newVersions, module.Version{
			Version: module.VersionNumber(tag),
		})
	}
	moduleMetadata.Versions = moduleMetadata.Versions.Merge(newVersions)

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
		OrganizationAddr: vcs.OrganizationAddr{
			Org: module.Namespace,
		},
		Name: "terraform-" + module.TargetSystem + "-" + module.Name,
	}
}
