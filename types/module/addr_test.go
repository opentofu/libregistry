// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package module_test

import (
	"encoding/json"
	"testing"

	"github.com/opentofu/libregistry/types/module"
)

func TestAddrJSON(t *testing.T) {
	const testNamespace = "opentofu"
	const testName = "test"
	const testSystem = "aws"

	moduleAddr := module.Addr{
		Namespace:    testNamespace,
		Name:         testName,
		TargetSystem: testSystem,
	}
	marshalled, err := json.Marshal(moduleAddr)
	if err != nil {
		t.Fatalf("Failed to marshal module address (%v)", err)
	}
	moduleAddr2 := module.Addr{}
	if err := json.Unmarshal(marshalled, &moduleAddr2); err != nil {
		t.Fatalf("Failed to unmarshal module address (%v)", err)
	}
	if !moduleAddr.Equals(moduleAddr2) {
		t.Fatalf("Module addresses are not equal.")
	}
}
