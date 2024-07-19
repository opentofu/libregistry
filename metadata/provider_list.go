// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) ListProviders(ctx context.Context) ([]provider.Addr, error) {
	providerLetters, err := r.storageAPI.ListDirectories(ctx, providersDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to list '%s' directory (%w)", providersDirectory, err)
	}
	var results []provider.Addr
	for _, letter := range providerLetters {
		letterResults, e := r.listProvidersByNamespaceLetter(ctx, letter)
		if e != nil {
			return nil, e
		}
		results = append(results, letterResults...)
	}
	return results, nil
}

func (r registryDataAPI) listProvidersByNamespaceLetter(ctx context.Context, letter string) ([]provider.Addr, error) {
	p := path.Join(providersDirectory, letter)
	namespaces, e := r.storageAPI.ListDirectories(ctx, storage.Path(p))
	if e != nil {
		return nil, fmt.Errorf("failed to list provider directory %s (%w)", p, e)
	}
	var results []provider.Addr
	for _, namespace := range namespaces {
		namespaceResults, e2 := r.ListProvidersByNamespace(ctx, namespace)
		if e2 != nil {
			return nil, e2
		}
		results = append(results, namespaceResults...)
	}
	return results, nil
}

func (r registryDataAPI) ListProvidersByNamespace(ctx context.Context, namespace string) ([]provider.Addr, error) {
	namespace = provider.NormalizeNamespace(namespace)

	p := path.Join(providersDirectory, namespace[0:1], namespace)
	files, err := r.storageAPI.ListFiles(ctx, storage.Path(p))
	if err != nil {
		return nil, fmt.Errorf("failed to list files in module name directory %s (%w)", p, err)
	}
	var result []provider.Addr
	for _, file := range files {
		if strings.HasSuffix(file, ".json") {
			result = append(result, provider.Addr{
				Namespace: namespace,
				Name:      strings.ToLower(strings.TrimSuffix(file, ".json")),
			})
		}
	}
	return result, nil
}
