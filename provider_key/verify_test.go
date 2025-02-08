// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"context"
	"testing"
	"time"

	"github.com/opentofu/libregistry/types/provider"
)

func TestProviderValidVerify(t *testing.T) {
	pkv := setupProviderCall(t, "/SHASumsURL/", "/SHASumsSignatureURL/")
	addr := provider.Addr{
		Name:      "test",
		Namespace: "opentofu",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, err := pkv.VerifyProvider(ctx, addr)
	if err != nil {
		t.Fatalf("Failed to verify provider: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("Data has the wrong size: %d", len(data))
	}

	if data[0].Version != "0.2.0" {
		t.Fatalf("Wrong version was returned %s", data[0])
	}
}

func TestProviderInvalidVerify(t *testing.T) {
	pkv := setupProviderCall(t, "/invalid/", "/invalid/")
	addr := provider.Addr{
		Name:      "test",
		Namespace: "opentofu",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, err := pkv.VerifyProvider(ctx, addr)
	if err != nil {
		t.Fatalf("Failed to verify provider: %v", err)
	}

	if len(data) != 0 {
		t.Fatalf("Data size should be 0, instead is: %d", len(data))
	}
}
