// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"testing"
)

func TestValidSignature(t *testing.T) {
	key := generateKey(t)
	pubKey := getPubKey(t, key)
	signature, data := generateSignedData(t, key, []byte("test\n"))

	pk, err := New(pubKey, nil)
	if err != nil {
		t.Fatalf("Failed to build ProviderKey (%v)", err)
	}

	err = pk.ValidateSignature(signature, data)
	if err != nil {
		t.Fatalf("Could not validate the signature (%v)", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	key := generateKey(t)
	pubKey := getPubKey(t, key)
	anotherKey := generateKey(t)
	signature, data := generateSignedData(t, anotherKey, []byte("invalid_signature\n"))

	pk, err := New(pubKey, nil)
	if err != nil {
		t.Fatalf("Failed to build ProviderKey (%v)", err)
	}

	err = pk.ValidateSignature(signature, data)
	if err == nil {
		t.Fatalf("Err should be non-nil (%v)", err)
	}
}
