// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"path"
	"strings"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) getProviderPath(providerAddr provider.Addr) storage.Path {
	return storage.Path(path.Join(modulesDirectory, providerAddr.Namespace[0:1], providerAddr.Name) + ".json")
}

func (r registryDataAPI) normalizeProviderNamespace(namespace string) string {
	return strings.ToLower(namespace)
}

func (r registryDataAPI) normalizeProviderName(name string) string {
	return strings.ToLower(name)
}

func (r registryDataAPI) normalizeProviderAddr(providerAddr provider.Addr) provider.Addr {
	return provider.Addr{
		Namespace: r.normalizeProviderNamespace(providerAddr.Namespace),
		Name:      r.normalizeProviderName(providerAddr.Name),
	}
}
