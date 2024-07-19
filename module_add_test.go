// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package libregistry_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/opentofu/libregistry"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/memory"
	"github.com/opentofu/libregistry/types/module"
	"github.com/opentofu/libregistry/vcs"
	"github.com/opentofu/libregistry/vcs/fakevcs"
	"github.com/opentofu/libregistry/vcs/github"
)

func ExampleAPI_AddModule() {
	ghClient, err := github.New(os.Getenv("GITHUB_TOKEN"), nil)
	if err != nil {
		panic(err)
	}

	storage := memory.New()

	dataAPI, err := metadata.New(storage)
	if err != nil {
		panic(err)
	}

	registry, err := libregistry.New(
		ghClient,
		dataAPI,
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	if err := registry.AddModule(ctx, "terraform-aws-modules/terraform-aws-iam"); err != nil {
		panic(err)
	}

	// Manually read the registry data file:
	jsonFile, err := storage.GetFile(ctx, "modules/t/terraform-aws-modules/iam/aws.json")
	if err != nil {
		panic(err)
	}

	var data map[string]any
	if err := json.Unmarshal(jsonFile, &data); err != nil {
		panic(err)
	}
	fmt.Printf("Latest version: %s", data["versions"].([]any)[0].(map[string]any)["version"].(string))
}

// TestAddModule tests that a module, when added from a repository, correctly appears in the metadata storage with
// the correct version information.
func TestAddModule(t *testing.T) {
	inMemoryVCS := fakevcs.New()
	storage := memory.New()
	ctx := context.Background()
	moduleAddr := module.Addr{
		Namespace:    "test",
		Name:         "aws",
		TargetSystem: "iam",
	}

	dataAPI, err := metadata.New(storage)
	if err != nil {
		panic(err)
	}

	registry, err := libregistry.New(
		inMemoryVCS,
		dataAPI,
	)
	if err != nil {
		t.Fatal(err)
	}

	org := vcs.OrganizationAddr{
		Org: moduleAddr.Namespace,
	}
	repo := vcs.RepositoryAddr{
		OrganizationAddr: org,
		Name:             "terraform-" + moduleAddr.TargetSystem + "-" + moduleAddr.Name,
	}

	if err := inMemoryVCS.CreateOrganization(org); err != nil {
		t.Fatal(err)
	}

	if err := inMemoryVCS.CreateRepository(repo); err != nil {
		t.Fatal(err)
	}

	if err := inMemoryVCS.CreateVersion(repo, "v1.0.0"); err != nil {
		t.Fatal(err)
	}

	if err := registry.AddModule(ctx, repo.String()); err != nil {
		t.Fatal(err)
	}

	storedMetadata, err := dataAPI.GetModule(ctx, moduleAddr)
	if err != nil {
		t.Fatal(err)
	}

	if len(storedMetadata.Versions) != 1 {
		t.Fatalf("Incorrect number of versions: %d", len(storedMetadata.Versions))
	}

	if ver := storedMetadata.Versions[0].Version; ver != "v1.0.0" {
		t.Fatalf("Incorrect version stored: %s", ver)
	}
}
