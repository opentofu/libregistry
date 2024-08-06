// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"encoding/json"
	"testing"

	"github.com/opentofu/libregistry/types/provider"
)

func TestAddrJSON(t *testing.T) {
	const testNamespace = "opentofu"
	const testName = "test"

	providerAddr := provider.Addr{
		Namespace: testNamespace,
		Name:      testName,
	}
	marshalled, err := json.Marshal(providerAddr)
	if err != nil {
		t.Fatalf("Failed to marshal provider address (%v)", err)
	}
	providerAddr2 := provider.Addr{}
	if err := json.Unmarshal(marshalled, &providerAddr2); err != nil {
		t.Fatalf("Failed to unmarshal provider address (%v)", err)
	}
	if !providerAddr.Equals(providerAddr2) {
		t.Fatalf("Module addresses are not equal.")
	}
}
