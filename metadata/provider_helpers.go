// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"path"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) getProviderPathRaw(providerAddr provider.Addr) storage.Path {
	providerAddr = providerAddr.Normalize()
	return storage.Path(path.Join(providersDirectory, providerAddr.Namespace[0:1], providerAddr.Namespace, providerAddr.Name) + ".json")
}

func (r registryDataAPI) getProviderPathCanonical(ctx context.Context, providerAddr provider.Addr) (storage.Path, error) {
	providerAddr, err := r.GetProviderCanonicalAddr(ctx, providerAddr)
	if err != nil {
		return "", err
	}
	return r.getProviderPathRaw(providerAddr), nil
}
