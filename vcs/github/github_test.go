// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package github_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/opentofu/libregistry/logger"
	"github.com/opentofu/libregistry/vcs"
	"github.com/opentofu/libregistry/vcs/github"
)

func TestClone(t *testing.T) {
	const testOrg = "integrations"
	const testRepo = "terraform-provider-github"
	const testVersion = "v6.2.3"

	checkoutDir := t.TempDir()

	gh, err := github.New(
		github.WithCheckoutRootDirectory(checkoutDir),
		github.WithLogger(logger.NewTestLogger(t)),
	)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	workingCopy, err := gh.Checkout(ctx, vcs.RepositoryAddr{
		Org:  testOrg,
		Name: testRepo,
	}, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = workingCopy.Close()
	})
	readme, err := workingCopy.Open("README.md")
	if err != nil {
		t.Fatal(err)
	}
	readmeContents, err := io.ReadAll(readme)
	if err != nil {
		t.Fatal(err)
	}
	if len(readmeContents) == 0 {
		t.Fatal("Empty readme!")
	}
}

func TestCloneNotFound(t *testing.T) {
	var notFound *vcs.RepositoryNotFoundError

	t.Logf("⚙️ Checking if cloning a non-existent repository correctly returns a %T...", notFound)

	const testOrg = "opentofu"
	const testRepo = "nonexistent"
	const testVersion = "v1.6.0"

	checkoutDir := t.TempDir()

	gh, err := github.New(
		github.WithCheckoutRootDirectory(checkoutDir),
		github.WithLogger(logger.NewTestLogger(t)),
	)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	_, err = gh.Checkout(ctx, vcs.RepositoryAddr{
		Org:  testOrg,
		Name: testRepo,
	}, testVersion)
	if err == nil {
		t.Fatal("❌ No error returned for nonexistent repository")
	}
	if !errors.As(err, &notFound) {
		t.Fatalf("❌ Cloning a non-existent repository did not return the correct error type (expected: %T, got: %T)", notFound, err)
	}
	t.Logf("✅ The cloning returned the correct error type for a non-existent repository.")
}
