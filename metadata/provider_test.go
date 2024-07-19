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
	const testNamespace = "opentofu"
	const testName = "test"

	fa := memory.New()
	api, err := metadata.New(fa)
	if err != nil {
		t.Fatalf("Failed to initialize API (%v)", err)
	}

	ctx := context.Background()

	checkEmpty := func(t *testing.T) {
		providerList, err := api.ListProviders(ctx)
		if err != nil {
			t.Fatalf("Failed to list providers: %v", err)
		}
		if len(providerList) != 0 {
			t.Fatalf("The provider list is not empty.")
		}

		providerList, err = api.ListProvidersByNamespace(ctx, testNamespace)
		if err != nil {
			t.Fatalf("Failed to list providers: %v", err)
		}
		if len(providerList) != 0 {
			t.Fatalf("The namespaced provider list is not empty.")
		}

		_, err = api.GetProvider(ctx, provider.Addr{
			Namespace: testNamespace,
			Name:      testName,
		})
		if err == nil {
			t.Fatalf("Getting a non-existent provider did not return an error.")
		}
		var typedErr *metadata.ProviderNotFoundError
		if !errors.As(err, &typedErr) {
			t.Fatalf("Fetching a non-existent provider did not return the correct error type (%T instead of %T)", err, typedErr)
		}
	}

	t.Run("1-list-get", checkEmpty)
}
