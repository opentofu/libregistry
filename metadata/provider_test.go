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

	fa := memory.New()
	api, err := metadata.New(fa)
	if err != nil {
		t.Fatalf("Failed to initialize API (%v)", err)
	}

	ctx := context.Background()

	checkEmpty := func(t *testing.T) {
		for _, ns := range []string{testAliasedNamespace, testNamespace} {
			providerList, err := api.ListProviders(ctx)
			if err != nil {
				t.Fatalf("Failed to list providers: %v", err)
			}
			if len(providerList) != 0 {
				t.Fatalf("The provider list is not empty.")
			}

			providerList, err = api.ListProvidersByNamespace(ctx, ns)
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
}
