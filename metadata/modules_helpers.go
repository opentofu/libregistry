// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"path"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/module"
)

func (r registryDataAPI) getModulePath(moduleAddr module.Addr) storage.Path {
	moduleAddr = moduleAddr.Normalize()
	return storage.Path(path.Join(modulesDirectory, moduleAddr.Namespace[0:1], moduleAddr.Namespace, moduleAddr.Name, moduleAddr.TargetSystem) + ".json")
}
