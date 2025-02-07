// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

type mockMetadata struct {
	metadata.API
}

func (m mockMetadata) GetProvider(ctx context.Context, addr provider.Addr, test bool) (provider.Metadata, error) {
	providerMetadata := provider.Metadata{
		Versions: provider.VersionList{
			provider.Version{
				Version: "0.2.0",
			},
		},
	}
	return providerMetadata, nil
}

func TestProviderVerify(t *testing.T) {
	metadataAPI := &mockMetadata{}
	httpClient := generateTestClient(t, "test")

	ctx := context.Background()

	pubKey := generateTestPubKey(t)

	pkv, err := New(pubKey, metadataAPI, WithHTTPClient(httpClient))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	addr := provider.Addr{
		Name:      "test",
		Namespace: "opentofu",
	}

	data, err := pkv.VerifyProvider(ctx, addr)
	if err != nil {
		t.Fatalf("Failed to verify key: %v", err)
	}

	if data[0].Version != "0.2.0" {
		t.Fatalf("Wrong version was returned %s", data[0])
	}
}
