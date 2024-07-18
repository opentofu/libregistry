package metadata_test

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/memory"
	"github.com/opentofu/libregistry/types/module"
)

func TestModuleCRUD(t *testing.T) {
	const testNamespace = "opentofu"
	const testName = "test"
	const testTargetSystem = "amd64"
	const testVersion = "1.0.0"

	fa := memory.New()
	api, err := metadata.New(fa)
	if err != nil {
		t.Fatalf("Failed to initialize API (%v)", err)
	}

	ctx := context.Background()

	checkEmpty := func(t *testing.T) {
		modules, err := api.ListModules(ctx)
		if err != nil {
			t.Fatalf("Failed to list modules (%v)", err)
		}
		if len(modules) != 0 {
			t.Fatalf("The module list is not empty.")
		}

		modules, err = api.ListModulesByNamespace(ctx, testNamespace)
		if err != nil {
			t.Fatalf("Failed to list modules (%v)", err)
		}
		if len(modules) != 0 {
			t.Fatalf("The module list is not empty.")
		}

		modules, err = api.ListModulesByNamespaceAndName(ctx, testNamespace, testName)
		if err != nil {
			t.Fatalf("Failed to list modules (%v)", err)
		}
		if len(modules) != 0 {
			t.Fatalf("The module list is not empty.")
		}

		_, err = api.GetModule(ctx, module.Addr{
			Namespace:    testNamespace,
			Name:         testName,
			TargetSystem: testTargetSystem,
		})
		if err == nil {
			t.Fatalf("Fetching a non-existent module did not result in an error.")
		}
	}
	t.Run("1-list-get", checkEmpty)
	t.Run("2-create", func(t *testing.T) {
		if err := api.PutModule(ctx, module.Addr{
			Namespace:    testNamespace,
			Name:         testName,
			TargetSystem: testTargetSystem,
		}, module.Metadata{
			Versions: []module.Version{
				{
					Version: testVersion,
				},
			},
		}); err != nil {
			t.Fatalf("Failed to put module (%v)", err)
		}
	})

	t.Run("3-list", func(t *testing.T) {
		checkModules := func(modules []module.Addr) {
			if len(modules) != 1 {
				t.Fatalf("The module list has %d elements.", len(modules))
			}
			if modules[0].Namespace != testNamespace {
				t.Fatalf("Invalid namespace: %s", modules[0].Namespace)
			}
			if modules[0].Name != testName {
				t.Fatalf("Invalid name: %s", modules[0].Name)
			}
			if modules[0].TargetSystem != testTargetSystem {
				t.Fatalf("Invalid target system: %s", modules[0].TargetSystem)
			}
		}

		modules, err := api.ListModules(ctx)
		if err != nil {
			t.Fatalf("Failed to list modules (%v)", err)
		}
		checkModules(modules)

		modules, err = api.ListModulesByNamespace(ctx, testNamespace)
		if err != nil {
			t.Fatalf("Failed to list modules (%v)", err)
		}
		checkModules(modules)
		modules, err = api.ListModulesByNamespace(ctx, testNamespace+"x")
		if err != nil {
			t.Fatalf("Failed to list modules (%v)", err)
		}
		if len(modules) != 0 {
			t.Fatalf("The module list is not empty.")
		}

		modules, err = api.ListModulesByNamespaceAndName(ctx, testNamespace, testName)
		if err != nil {
			t.Fatalf("Failed to list modules (%v)", err)
		}
		checkModules(modules)

		modules, err = api.ListModulesByNamespaceAndName(ctx, testNamespace, testName+"x")
		if err != nil {
			t.Fatalf("Failed to list modules (%v)", err)
		}
		if len(modules) != 0 {
			t.Fatalf("The module list is not empty.")
		}
	})
	t.Run("4-get", func(t *testing.T) {
		mod, err := api.GetModule(ctx, module.Addr{
			Namespace:    testNamespace,
			Name:         testName,
			TargetSystem: testTargetSystem,
		})
		if err != nil {
			t.Fatalf("Failed to get module (%v)", err)
		}
		if len(mod.Versions) != 1 {
			t.Fatalf("Incorrect number of module versions: %d", len(mod.Versions))
		}
		if mod.Versions[0].Version != testVersion {
			t.Fatalf("Incorrect module version: %s", mod.Versions[0].Version)
		}

		if _, err = api.GetModule(ctx, module.Addr{
			Namespace:    testNamespace + "x",
			Name:         testName,
			TargetSystem: testTargetSystem,
		}); err == nil {
			t.Fatalf("Getting invalid module namespace did not result in an error.")
		}
		if _, err = api.GetModule(ctx, module.Addr{
			Namespace:    testNamespace,
			Name:         testName + "x",
			TargetSystem: testTargetSystem,
		}); err == nil {
			t.Fatalf("Getting invalid module name did not result in an error.")
		}
		if _, err = api.GetModule(ctx, module.Addr{
			Namespace:    testNamespace,
			Name:         testName,
			TargetSystem: testTargetSystem + "x",
		}); err == nil {
			t.Fatalf("Getting invalid module target system did not result in an error.")
		}
	})
	t.Run("5-delete", func(t *testing.T) {
		if err := api.DeleteModule(ctx, module.Addr{
			Namespace:    testNamespace,
			Name:         testName,
			TargetSystem: testTargetSystem,
		}); err != nil {
			t.Fatalf("Failed to delete module (%v)", err)
		}

		if err := api.DeleteModule(ctx, module.Addr{
			Namespace:    testNamespace,
			Name:         testName,
			TargetSystem: testTargetSystem,
		}); err != nil {
			t.Fatalf("Failed to delete already-deleted module (%v)", err)
		}
	})
	t.Run("6-list-get", checkEmpty)
}
