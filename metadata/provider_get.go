// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

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
		_, err = r.storageAPI.FileExists(ctx, providerPath)
		if err != nil {
			return providerAddr, providerPath, err
		}
	}

	if targetAddr, ok := providerAliases[providerAddr]; ok {
		providerAddr = targetAddr.Normalize()
		providerPath = r.getProviderPathRaw(providerAddr)
		_, err = r.storageAPI.FileExists(ctx, providerPath)
		if err != nil {
			return providerAddr, providerPath, err
		}
	}

	exists, err := r.storageAPI.FileExists(ctx, providerPath)
	if err != nil {
		return providerAddr, providerPath, err
	}
	if exists {
		return providerAddr, providerPath, nil
	}

	return providerAddr, providerPath, &ProviderNotFoundError{
		ProviderAddr: providerAddr,
	}
}

func (r registryDataAPI) GetProviderReverseAliases(ctx context.Context, addr provider.Addr) ([]provider.Addr, error) {
	providerAliases, err := r.ListProviderAliases(ctx)
	if err != nil {
		return nil, err
	}
	namespaceAliases, err := r.ListProviderNamespaceAliases(ctx)
	if err != nil {
		return nil, err
	}

	var results []provider.Addr
	for alias, target := range providerAliases {
		if target.Equals(addr) {
			results = append(results, alias)
			// Look up the namespace alias of the provider alias to include it in the list.
			for namespace, targetNamespace := range namespaceAliases {
				if provider.NormalizeNamespace(targetNamespace) == provider.NormalizeNamespace(alias.Namespace) {
					results = append(results, provider.Addr{
						Namespace: namespace,
						Name:      target.Name,
					})
				}
			}

		}
	}
	for namespace, targetNamespace := range namespaceAliases {
		if provider.NormalizeNamespace(targetNamespace) == provider.NormalizeNamespace(addr.Namespace) {
			results = append(results, provider.Addr{
				Namespace: namespace,
				Name:      addr.Name,
			})
		}
	}
	return results, nil
}
