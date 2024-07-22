// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

const maxProviderAliasRecursionDepth = 5

func (r registryDataAPI) GetProvider(ctx context.Context, providerAddr provider.Addr, resolveAliases bool) (provider.Metadata, error) {
	var path storage.Path
	var err error
	if resolveAliases {
		path, err = r.getProviderPathCanonical(ctx, providerAddr)
		if err != nil {
			return provider.Metadata{}, err
		}
	} else {
		path = r.getProviderPathRaw(providerAddr)
	}

	fileContents, err := r.storageAPI.GetFile(ctx, path)
	if err != nil {
		var notFound *storage.ErrFileNotFound
		if errors.As(err, &notFound) {
			return provider.Metadata{}, &ProviderNotFoundError{
				ProviderAddr: providerAddr,
				Cause:        err,
			}
		}
		return provider.Metadata{}, fmt.Errorf("failed to read provider file %s (%w)", path, err)
	}
	var mod provider.Metadata
	if err := json.Unmarshal(fileContents, &mod); err != nil {
		return provider.Metadata{}, fmt.Errorf("failed to parse provider metadata file %s (%w)", path, err)
	}
	return mod, nil
}

func (r registryDataAPI) GetProviderCanonicalAddr(ctx context.Context, providerAddr provider.Addr) (provider.Addr, error) {
	providerAddr, _, err := r.getProviderCanonical(ctx, providerAddr)
	return providerAddr, err
}

func (r registryDataAPI) getProviderCanonical(ctx context.Context, providerAddr provider.Addr) (provider.Addr, storage.Path, error) {
	providerAddr = providerAddr.Normalize()
	providerPath := r.getProviderPathRaw(providerAddr)
	exists, err := r.storageAPI.FileExists(ctx, providerPath)
	if err != nil {
		return providerAddr, providerPath, err
	}
	if exists {
		return providerAddr, providerPath, nil
	}

	namespaceAliases, err := r.ListProviderNamespaceAliases(ctx)
	if err != nil {
		return providerAddr, providerPath, err
	}

	providerAliases, err := r.ListProviderAliases(ctx)
	if err != nil {
		return providerAddr, providerPath, err
	}

	providerAddr = providerAddr.Normalize()

	if targetNamespace, ok := namespaceAliases[providerAddr.Namespace]; ok {
		providerAddr = provider.Addr{
			Namespace: targetNamespace,
			Name:      providerAddr.Name,
		}.Normalize()
		providerPath = r.getProviderPathRaw(providerAddr)
		exists, err = r.storageAPI.FileExists(ctx, providerPath)
		if err != nil {
			return providerAddr, providerPath, err
		}
		if exists {
			// If the aliased provider exists, don't look any further.
			return providerAddr, providerPath, nil
		}
	}

	if targetAddr, ok := providerAliases[providerAddr]; ok {
		providerAddr = targetAddr.Normalize()
		providerPath = r.getProviderPathRaw(providerAddr)
		exists, err = r.storageAPI.FileExists(ctx, providerPath)
		if err != nil {
			return providerAddr, providerPath, err
		}
		if exists {
			// If the aliased provider exists, don't look any further.
			return providerAddr, providerPath, nil
		}
	}

	return providerAddr, providerPath, &ProviderNotFoundError{
		ProviderAddr: providerAddr,
	}
}
