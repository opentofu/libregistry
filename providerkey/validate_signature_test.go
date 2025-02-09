// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/types/provider"
)

func TestValidSignature(t *testing.T) {
	key := generateKey(t)
	pubKey := getPubKey(t, key)
	signature, data := generateSignedData(t, key, []byte("test\n"))

	pk, err := New(pubKey, nil)
	if err != nil {
		t.Fatalf("Failed to build ProviderKey (%v)", err)
	}

	p := provider.Addr{Name: "test"}
	err = pk.ValidateSignature(context.Background(), p, signature, data)
	if err != nil {
		t.Fatalf("Could not validate the signature (%v)", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	key1 := generateKey(t)
	signature, data := generateSignedData(t, key1, []byte("test\n"))

	key2 := generateKey(t)
	pubKey2 := getPubKey(t, key2)
	pk, err := New(pubKey2, nil)
	if err != nil {
		t.Fatalf("Failed to build ProviderKey (%v)", err)
	}

	p := provider.Addr{Name: "test"}
	err = pk.ValidateSignature(context.Background(), p, signature, data)
	if err == nil {
		t.Fatalf("Err should be non-nil (%v)", err)
	}
}
