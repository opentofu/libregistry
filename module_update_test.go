// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package libregistry_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/opentofu/libregistry"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/memory"
	"github.com/opentofu/libregistry/types/module"
	"github.com/opentofu/libregistry/vcs"
	"github.com/opentofu/libregistry/vcs/fakevcs"
)

// TestUpdateModuleBackfill tests that the process correctly backfills using the expensive ListAllVersions call when
// it detects that there was no overlap between the existing versions and the versions returned by ListLatestVersions.
// This test relies on the fact that the fakevcs returns 5 versions in the ListLatestVersions call.
func TestUpdateModuleBackfill(t *testing.T) {
	const createVersionCount = 20

	moduleAddr := module.Addr{
		Namespace:    "test",
		Name:         "aws",
		TargetSystem: "iam",
	}
	org := vcs.OrganizationAddr{
		Org: moduleAddr.Namespace,
	}
	repo := vcs.RepositoryAddr{
		OrganizationAddr: org,
		Name:             "terraform-" + moduleAddr.TargetSystem + "-" + moduleAddr.Name,
	}

	inMemoryVCS := fakevcs.New()
	storage := memory.New()
	ctx := context.Background()
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

	for i := 1; i <= createVersionCount; i++ {
		if err := inMemoryVCS.CreateVersion(repo, "v1.0."+strconv.Itoa(i)); err != nil {
			t.Fatal(err)
		}
	}

	if err := registry.UpdateModule(ctx, moduleAddr); err != nil {
		t.Fatal(err)
	}

	storedMetadata, err := dataAPI.GetModule(ctx, moduleAddr)
	if err != nil {
		t.Fatal(err)
	}

	if len(storedMetadata.Versions) != createVersionCount+1 {
		t.Fatalf("Incorrect number of versions: %d", len(storedMetadata.Versions))
	}

	j := 0
	for i := createVersionCount; i >= 0; i-- {
		if ver := storedMetadata.Versions[j].Version; ver != module.VersionNumber("v1.0."+strconv.Itoa(i)) {
			t.Fatalf("Incorrect version in position %d: %s", j, ver)
		}
		j++
	}
}