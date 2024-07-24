// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata_test

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/memory"
	"github.com/opentofu/libregistry/types/provider"
)

func TestProviderAPI_GetAllProviders(t *testing.T) {
	const testNamespace = "opentofu"
	const testName = "test"

	providerAddr := provider.Addr{
		Namespace: testNamespace,
		Name:      testName,
	}
	providerVersion := provider.Version{
		Version:             "v1.0.0",
		Protocols:           []string{"5.0"},
		SHASumsURL:          "https://localhost/" + providerAddr.Namespace + "/" + providerAddr.Name + "/releases/download/v1.0.0/" + providerAddr.String() + "_SHA256SUMS",
		SHASumsSignatureURL: "https://localhost/" + providerAddr.Namespace + "/" + providerAddr.Name + "/releases/download/v1.0.0/" + providerAddr.String() + "_SHA256SUMS.sig",
		Targets: []provider.Target{
			{
				OS:          "linux",
				Arch:        "amd64",
				Filename:    providerAddr.String() + "_linux_amd64.zip",
				DownloadURL: "https://localhost/" + providerAddr.Namespace + "/" + providerAddr.Name + "/releases/download/v1.0.0/" + providerAddr.String() + "_linux_amd64.zip",
				SHASum:      "c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a",
			},
		},
	}
	providerMetadata := provider.Metadata{
		CustomRepository: "",
		Versions: []provider.Version{
			providerVersion,
		},
	}

	storage := memory.New()
	api, err := metadata.New(storage)
	if err != nil {
		t.Fatalf("Failed to initialize API (%v)", err)
	}

	ctx := context.Background()

	if err := api.PutProvider(ctx, providerAddr, providerMetadata); err != nil {
		t.Fatalf("Failed to put provider (%v)", err)
	}

	allProviders, err := api.GetAllProviders(ctx, true)
	if err != nil {
		t.Fatalf("Failed to get all providers (%v)", err)
	}
	if n := len(allProviders); n != 1 {
		t.Fatalf("Incorrect number of providers in registry (%d)", n)
	}

	foundMetadata := allProviders[providerAddr]
	if !foundMetadata.Equals(providerMetadata) {
		t.Fatalf("Incorrect provider metadata returned.")
	}
}
