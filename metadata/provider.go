// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"

	"github.com/opentofu/libregistry/types/provider"
)

// ProviderDataAPI lists the methods for handling providers and their keys.
type ProviderDataAPI interface {
	// ListProviderNamespaceAliases returns a list of source to target namespace aliases. The key is the "from"
	// namespace, the value is the "to" namespace. The alias means that all providers in the "to" namespace should
	// also be observed in the "from" namespace.
	ListProviderNamespaceAliases(ctx context.Context) (map[string]string, error)

	// ListProviders returns all providers in the registry. The includeAliases parameter lets you include aliased copies
	// of providers.
	ListProviders(ctx context.Context, includeAliases bool) ([]provider.Addr, error)
	// ListProvidersByNamespace returns all providers in a given namespace, returning the addresses. The includeAliases
	// parameter lets you include aliased copies of providers.
	ListProvidersByNamespace(ctx context.Context, namespace string, includeAliases bool) ([]provider.Addr, error)

	// GetProvider returns the metadata for a given provider address. The resolveAliases parameter lets you control
	// if provider namespace aliases should be resolved or not.
	GetProvider(ctx context.Context, addr provider.Addr, resolveAliases bool) (provider.Metadata, error)
	// GetProviderCanonicalAddr returns the canonical address of a provider, resolving any namespace aliases and
	// lowercasing the name. This function may return a *ProviderNotFoundError if the provider was not found in
	// its original or in the target namespace.
	GetProviderCanonicalAddr(ctx context.Context, addr provider.Addr) (provider.Addr, error)

	// PutProvider queues up writing the specified provider metadata.
	PutProvider(ctx context.Context, addr provider.Addr, metadata provider.Metadata) error
	// DeleteProvider queues up deleting the specified provider.
	DeleteProvider(ctx context.Context, addr provider.Addr) error

	// ListProviderNamespacesWithKeys returns a list of provider namespaces that have a key registered.
	ListProviderNamespacesWithKeys(ctx context.Context) ([]string, error)
	// ListProviderNamespaceKeyIDs lists the keys IDs of all keys registered in a provider namespace.
	ListProviderNamespaceKeyIDs(ctx context.Context, namespace string) ([]string, error)

	// GetProviderNamespaceKey loads a key for a specific provider namespace and returns the key material.
	GetProviderNamespaceKey(ctx context.Context, namespace string, keyID string) (provider.Key, error)

	// PutProviderNamespaceKey queues up adding a key with the specified key material for a provider namespace.
	PutProviderNamespaceKey(ctx context.Context, namespace string, key provider.Key) error
	// DeleteProviderNamespaceKey queues up deleting a specific key from a provider namespace.
	DeleteProviderNamespaceKey(ctx context.Context, namespace string, keyID string) error
}

const providersDirectory = "providers"
