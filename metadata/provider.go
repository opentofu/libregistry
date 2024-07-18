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
	// namespace, the value is the "to" namespace.
	ListProviderNamespaceAliases(ctx context.Context) (map[string]string, error)
	// PutProviderNamespaceAlias queues up adding an alias. The "to" namespace will be used for storage.
	// Any providers registered in the "from" namespace will be moved to the "to" namespace.
	PutProviderNamespaceAlias(ctx context.Context, from string, to string) error
	// DeleteProviderNamespaceAlias queues up deleting an alias from the specified "from" namespace.
	DeleteProviderNamespaceAlias(ctx context.Context, from string) error

	// ListProviders returns all providers in the registry.
	ListProviders(ctx context.Context) ([]provider.Addr, error)
	// ListProvidersByNamespace returns all providers in a given namespace, returning the addresses.
	ListProvidersByNamespace(ctx context.Context, namespace string) ([]provider.Addr, error)

	// GetProvider returns the metadata for a given provider address.
	GetProvider(ctx context.Context, addr provider.Addr) (provider.Metadata, error)
	// GetProviderCanonicalAddr returns the canonical address of a provider, resolving any namespace aliases and
	// lowercasing the name.
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