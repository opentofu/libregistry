// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) ListProviders(ctx context.Context, includeAliases bool) ([]provider.Addr, error) {
	providerLetters, err := r.storageAPI.ListDirectories(ctx, providersDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to list '%s' directory (%w)", providersDirectory, err)
	}

	var namespaceAliases map[string]string
	var providerAliases map[provider.Addr]provider.Addr

	if includeAliases {
		// Make sure we have the aliases in our letter list.
		providerLetterSet := map[string]struct{}{}
		for _, letter := range providerLetters {
			providerLetterSet[letter] = struct{}{}
		}

		namespaceAliases, err = r.ListProviderNamespaceAliases(ctx)
		if err != nil {
			return nil, err
		}
		for from := range namespaceAliases {
			providerLetterSet[from[:1]] = struct{}{}
		}

		providerAliases, err = r.ListProviderAliases(ctx)
		if err != nil {
			return nil, err
		}
		for alias := range providerAliases {
			providerLetterSet[alias.Namespace[:1]] = struct{}{}
		}

		providerLetters = make([]string, len(providerLetterSet))
		i := 0
		for letter := range providerLetterSet {
			providerLetters[i] = letter
			i++
		}
	}

	var results []provider.Addr
	for _, letter := range providerLetters {
		letterResults, e := r.listProvidersByNamespaceLetter(ctx, letter, includeAliases, namespaceAliases, providerAliases)
		if e != nil {
			return nil, e
		}
		results = append(results, letterResults...)
	}
	return results, nil
}

func (r registryDataAPI) listProvidersByNamespaceLetter(
	ctx context.Context,
	letter string,
	includeAliases bool,
	namespaceAliases map[string]string,
	providerAliases map[provider.Addr]provider.Addr,
) ([]provider.Addr, error) {
	p := path.Join(providersDirectory, letter)
	namespaces, e := r.storageAPI.ListDirectories(ctx, storage.Path(p))
	if e != nil {
		return nil, fmt.Errorf("failed to list provider directory %s (%w)", p, e)
	}

	if includeAliases {
		// Make sure we have the aliases in our letter list.
		providerNamespaceSet := map[string]struct{}{}
		for _, namespace := range namespaces {
			providerNamespaceSet[namespace] = struct{}{}
		}

		for from := range namespaceAliases {
			if from[:1] == letter {
				providerNamespaceSet[from] = struct{}{}
			}
		}
		for from := range providerAliases {
			if from.Namespace[:1] == letter {
				providerNamespaceSet[from.Namespace] = struct{}{}
			}
		}

		namespaces = make([]string, len(providerNamespaceSet))
		i := 0
		for namespace := range providerNamespaceSet {
			namespaces[i] = namespace
			i++
		}
	}

	var results []provider.Addr
	for _, namespace := range namespaces {
		namespaceResults, e2 := r.listProvidersByNamespace(ctx, namespace, includeAliases, includeAliases, namespaceAliases, providerAliases)
		if e2 != nil {
			return nil, e2
		}
		results = append(results, namespaceResults...)
	}
	return results, nil
}

func (r registryDataAPI) ListProvidersByNamespace(ctx context.Context, namespace string, includeAliases bool) ([]provider.Addr, error) {
	var err error
	var namespaceAliases map[string]string
	var providerAliases map[provider.Addr]provider.Addr

	if includeAliases {
		namespaceAliases, err = r.ListProviderNamespaceAliases(ctx)
		if err != nil {
			return nil, err
		}

		providerAliases, err = r.ListProviderAliases(ctx)
		if err != nil {
			return nil, err
		}
	}
	return r.listProvidersByNamespace(ctx, namespace, includeAliases, includeAliases, namespaceAliases, providerAliases)
}

func (r registryDataAPI) listProvidersByNamespace(
	ctx context.Context,
	namespace string,
	includeNamespaceAliases bool,
	includeProviderAliases bool,
	namespaceAliases map[string]string,
	providerAliases map[provider.Addr]provider.Addr,
) ([]provider.Addr, error) {
	namespace = provider.NormalizeNamespace(namespace)

	p := path.Join(providersDirectory, namespace[0:1], namespace)
	files, err := r.storageAPI.ListFiles(ctx, storage.Path(p))
	if err != nil {
		return nil, fmt.Errorf("failed to list files in module name directory %s (%w)", p, err)
	}
	var result []provider.Addr
	providerAddrSet := map[provider.Addr]struct{}{}
	for _, file := range files {
		if strings.HasSuffix(file, ".json") {
			addr := provider.Addr{
				Namespace: namespace,
				Name:      strings.ToLower(strings.TrimSuffix(file, ".json")),
			}
			result = append(result, addr)
			providerAddrSet[addr] = struct{}{}
		}
	}

	if includeNamespaceAliases {
		if target, ok := namespaceAliases[namespace]; ok {
			aliasedProviderAddrs, err := r.listProvidersByNamespace(ctx, target, false, true, namespaceAliases, providerAliases)
			if err != nil {
				return nil, err
			}
			for _, addr := range aliasedProviderAddrs {
				addr.Namespace = namespace
				if _, present := providerAddrSet[addr]; !present {
					result = append(result, addr)
					providerAddrSet[addr] = struct{}{}
				}
			}
		}
	}

	if includeProviderAliases {
		for from, to := range providerAliases {
			if from.Namespace == namespace {
				addr := provider.Addr{
					Namespace: namespace,
					Name:      to.Name,
				}
				_, err = r.getProviderPathCanonical(ctx, addr)
				if err != nil {
					var notFound *ProviderNotFoundError
					if !errors.As(err, &notFound) {
						return nil, err
					}
				} else {
					if _, present := providerAddrSet[addr]; !present {
						result = append(result, addr)
						providerAddrSet[addr] = struct{}{}
					}
				}
			}
		}
	}

	return result, nil
}
