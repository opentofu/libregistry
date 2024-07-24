// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata_test

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/memory"
	"github.com/opentofu/libregistry/types/module"
)

func TestModuleAPI_GetAllModules(t *testing.T) {
	const testNamespace = "opentofu"
	const testName = "test"
	const testTargetSystem = "amd64"
	const testVersion = "1.0.0"

	moduleAddr := module.Addr{
		Namespace:    testNamespace,
		Name:         testName,
		TargetSystem: testTargetSystem,
	}
	metadataEntry := module.Metadata{
		Versions: []module.Version{
			module.Version{
				Version: testVersion,
			},
		},
	}

	storage := memory.New()
	api, err := metadata.New(storage)
	if err != nil {
		t.Fatalf("Failed to initialize API (%v)", err)
	}

	ctx := context.Background()

	if err := api.PutModule(ctx, moduleAddr, metadataEntry); err != nil {
		t.Fatalf("Failed to put module (%v)", err)
	}

	allModules, err := api.GetAllModules(ctx)
	if err != nil {
		t.Fatalf("Failed to get all modules (%v)", err)
	}

	if n := len(allModules); n != 1 {
		t.Fatalf("Incorrect number of modules in the registry: %d", n)
	}

	foundMetadata := allModules[moduleAddr]
	if !foundMetadata.Equals(metadataEntry) {
		t.Fatalf("Metadata entries are not equal!")
	}
}
