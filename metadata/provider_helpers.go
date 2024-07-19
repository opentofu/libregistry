// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"path"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) getProviderPath(providerAddr provider.Addr) storage.Path {
	providerAddr = providerAddr.Normalize()
	return storage.Path(path.Join(modulesDirectory, providerAddr.Namespace[0:1], providerAddr.Name) + ".json")
}
