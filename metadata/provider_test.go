// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata_test

import (
	"context"
	"errors"
	"testing"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/memory"
	"github.com/opentofu/libregistry/types/provider"
)

func TestProviderCRUD(t *testing.T) {
	const testAliasedNamespace = "opentofu"
	const testNamespace = "hashicorp"
	const testName = "test"

	canonicalAddr := provider.Addr{
		Namespace: testNamespace,
		Name:      testName,
	}
	aliasedAddr := provider.Addr{
		Namespace: testAliasedNamespace,
		Name:      testName,
	}
	providerVersion := provider.Version{
		Version:             "v1.0.0",
		Protocols:           []string{"5.0"},
		SHASumsURL:          "https://localhost/" + canonicalAddr.Namespace + "/" + canonicalAddr.Name + "/releases/download/v1.0.0/" + canonicalAddr.String() + "_SHA256SUMS",
		SHASumsSignatureURL: "https://localhost/" + canonicalAddr.Namespace + "/" + canonicalAddr.Name + "/releases/download/v1.0.0/" + canonicalAddr.String() + "_SHA256SUMS.sig",
		Targets: []provider.Target{
			{
				OS:          "linux",
				Arch:        "amd64",
				Filename:    canonicalAddr.String() + "_linux_amd64.zip",
				DownloadURL: "https://localhost/" + canonicalAddr.Namespace + "/" + canonicalAddr.Name + "/releases/download/v1.0.0/" + canonicalAddr.String() + "_linux_amd64.zip",
				SHASum:      "c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a",
			},
		},
	}

	storage := memory.New()
	api, err := metadata.New(storage)
	if err != nil {
		t.Fatalf("Failed to initialize API (%v)", err)
	}

	ctx := context.Background()

	checkEmpty := func(t *testing.T) {
		for _, ns := range []string{testAliasedNamespace, testNamespace} {
			providerList, err := api.ListProviders(ctx, false)
			if err != nil {
				t.Fatalf("Failed to list providers: %v", err)
			}
			if len(providerList) != 0 {
				t.Fatalf("The provider list is not empty.")
			}

			providerList, err = api.ListProvidersByNamespace(ctx, ns, false)
			if err != nil {
				t.Fatalf("Failed to list providers: %v", err)
			}
			if len(providerList) != 0 {
				t.Fatalf("The namespaced provider list is not empty.")
			}

			_, err = api.GetProvider(ctx, provider.Addr{
				Namespace: ns,
				Name:      testName,
			}, false)
			if err == nil {
				t.Fatalf("Getting a non-existent provider did not return an error.")
			}
			var typedErr *metadata.ProviderNotFoundError
			if !errors.As(err, &typedErr) {
				t.Fatalf("Fetching a non-existent provider did not return the correct error type (%T instead of %T)", err, typedErr)
			}
		}

		_, err := api.GetProviderCanonicalAddr(ctx, aliasedAddr)
		if err == nil {
			t.Fatalf("Canonical address resolution did not reutnr an error.")
		}
		var typedErr *metadata.ProviderNotFoundError
		if !errors.As(err, &typedErr) {
			t.Fatalf("Fetching the provider alias for a non-existent provider did not return the correct error type (%T instead of %T)", err, typedErr)
		}
	}

	t.Run("1-list-get", checkEmpty)
	t.Run("2-create", func(t *testing.T) {
		if err := api.PutProvider(ctx, canonicalAddr, provider.Metadata{
			CustomRepository: "",
			Versions: []provider.Version{
				providerVersion,
			},
		}); err != nil {
			t.Fatalf("Failed to create provider version (%v)", err)
		}
	})
	t.Run("3-list-get", func(t *testing.T) {
		providers, err := api.ListProviders(ctx, false)
		if err != nil {
			t.Fatalf("Failed to list providers (%v)", err)
		}
		if len(providers) != 1 {
			t.Fatalf("Incorrect number of providers in the registry (%d)", len(providers))
		}
		if !providers[0].Equals(canonicalAddr) {
			t.Fatalf("Incorrect provider addr in the registry (%s instead of %s)", providers[0].String(), canonicalAddr.String())
		}

		providers, err = api.ListProvidersByNamespace(ctx, testNamespace, false)
		if err != nil {
			t.Fatalf("Failed to list providers (%v)", err)
		}
		if len(providers) != 1 {
			t.Fatalf("Incorrect number of providers in the registry (%d)", len(providers))
		}
		if !providers[0].Equals(canonicalAddr) {
			t.Fatalf("Incorrect provider addr in the registry (%s instead of %s)", providers[0].String(), canonicalAddr.String())
		}

		providerMeta, err := api.GetProvider(ctx, canonicalAddr, false)
		if err != nil {
			t.Fatalf("Failed to get provider (%v)", err)
		}
		if !providerMeta.Equals(provider.Metadata{
			CustomRepository: "",
			Versions: []provider.Version{
				providerVersion,
			},
		}) {
			t.Fatalf("Incorrect provider metadata returned.")
		}

		_, err = api.GetProvider(ctx, aliasedAddr, false)
		if err == nil {
			t.Fatalf("No error returned when querying a provider by its alias without resolveAlias.")
		}
		var notFound *metadata.ProviderNotFoundError
		if !errors.As(err, &notFound) {
			t.Fatalf("Incorrect error type returned when querying a provider by its alias without resolveAlias (%T instead of %T)", err, notFound)
		}

		providers, err = api.ListProviders(ctx, true)
		if err != nil {
			t.Fatalf("Failed to list providers (%v)", err)
		}
		if len(providers) != 2 {
			t.Fatalf("Incorrect number of providers in the registry (%d)", len(providers))
		}
		if !providers[0].Equals(canonicalAddr) && !providers[1].Equals(canonicalAddr) {
			t.Fatalf("None of the returned addresses contained the canonical address.")
		}
		if !providers[0].Equals(aliasedAddr) && !providers[1].Equals(aliasedAddr) {
			t.Fatalf("None of the returned addresses contained the aliased address.")
		}

		providers, err = api.ListProvidersByNamespace(ctx, testNamespace, false)
		if err != nil {
			t.Fatalf("Failed to list providers (%v)", err)
		}
		if len(providers) != 1 {
			t.Fatalf("Incorrect number of providers in the registry (%d)", len(providers))
		}
		if !providers[0].Equals(canonicalAddr) {
			t.Fatalf("Incorrect provider addr in the registry (%s instead of %s)", providers[0].String(), canonicalAddr.String())
		}
		providers, err = api.ListProvidersByNamespace(ctx, testAliasedNamespace, true)
		if err != nil {
			t.Fatalf("Failed to list providers (%v)", err)
		}
		if len(providers) != 1 {
			t.Fatalf("Incorrect number of providers in the registry (%d)", len(providers))
		}
		if !providers[0].Equals(aliasedAddr) {
			t.Fatalf("Incorrect provider addr in the registry (%s instead of %s)", providers[0].String(), aliasedAddr.String())
		}

		aliasedProviderMeta, err := api.GetProvider(ctx, aliasedAddr, true)
		if err != nil {
			t.Fatalf("Failed to query the provider by its alias (%v)", err)
		}
		if !providerMeta.Equals(aliasedProviderMeta) {
			t.Fatalf("Aliased provider meta lookup returned different metadata.")
		}
	})
	t.Run("4-delete", func(t *testing.T) {
		if err := api.DeleteProvider(ctx, canonicalAddr); err != nil {
			t.Fatalf("Failed to delete provider (%v)", err)
		}
		if err := api.DeleteProvider(ctx, canonicalAddr); err != nil {
			t.Fatalf("Deleting an already-deleted provider failed (%v)", err)
		}
	})
	t.Run("5-list-get", checkEmpty)
}
