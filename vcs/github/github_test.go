// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package github_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/opentofu/libregistry/logger"
	"github.com/opentofu/libregistry/vcs"
	"github.com/opentofu/libregistry/vcs/github"
)

func TestRepoInfo(t *testing.T) {
	t.Parallel()

	const testOrg = "opentofu"
	const testRepo = "opentofu"

	t.Logf("⚙️ Checking if the GitHub API returns a repository description for %s/%s...", testOrg, testRepo)

	gh, err := github.New(
		github.WithLogger(logger.NewTestLogger(t)),
		github.WithToken(os.Getenv("GITHUB_TOKEN")),
	)
	if err != nil {
		t.Fatalf("❌ Failed to initialize Github client (%v)", err)
	}
	ctx := context.Background()

	info, err := gh.GetRepositoryInfo(ctx, vcs.RepositoryAddr{
		Org:  testOrg,
		Name: testRepo,
	})
	if err != nil {
		t.Fatalf("❌ Failed to fetch repository info for %s/%s (%v)", testOrg, testRepo, err)
	}
	if info.Description == "" {
		t.Fatalf("❌ No description returned for %s/%s.", testOrg, testRepo)
	}
	t.Logf("✅ The GitHub API returned the following description for %s/%s: %s", testOrg, testRepo, info.Description)
}

func TestReleases(t *testing.T) {
	const testOrg = "integrations"
	const testRepo = "terraform-provider-github"
	const testVersion = "v6.2.3"
	testDate, err := time.Parse(time.RFC3339, "2024-07-08T16:58:50Z")
	if err != nil {
		t.Fatalf("❌ Failed to parse test date (%v)", err)
	}

	t.Logf("⚙️ Checking if version %s is present in %s/%s and was released on %s...", testVersion, testOrg, testRepo, testDate.String())

	gh, err := github.New(
		github.WithLogger(logger.NewTestLogger(t)),
		github.WithToken(os.Getenv("GITHUB_TOKEN")),
	)
	if err != nil {
		t.Fatalf("❌ Failed to initialize Github client (%v)", err)
	}
	ctx := context.Background()

	releases, err := gh.ListAllReleases(ctx, vcs.RepositoryAddr{
		Org:  testOrg,
		Name: testRepo,
	})
	if err != nil {
		t.Fatalf("❌ Failed to list GitHub releases (%v)", err)
	}
	found := false
	for _, release := range releases {
		if release.VersionNumber.Equals(testVersion) {
			found = true
			if !release.Created.Equal(testDate) {
				t.Fatalf(
					"❌ Found version %s, but the release date is incorrect: %s (expected: %s)",
					testVersion,
					release.Created.String(),
					testDate.String(),
				)
			}
		}
	}
	if !found {
		t.Fatalf("❌ Expected version not found (%s)", testVersion)
	}
	t.Logf("✅ Found version %s with release date %s.", testVersion, testDate.String())
}

func TestTags(t *testing.T) {
	t.Parallel()

	const testOrg = "integrations"
	const testRepo = "terraform-provider-github"
	const testVersion = "v6.2.3"
	testDate, err := time.Parse(time.RFC3339, "2024-07-08T16:55:36Z")
	if err != nil {
		t.Fatalf("❌ Failed to parse test date (%v)", err)
	}

	t.Logf("⚙️ Checking if version %s is present in %s/%s and was released on %s...", testVersion, testOrg, testRepo, testDate.String())

	gh, err := github.New(
		github.WithLogger(logger.NewTestLogger(t)),
		github.WithToken(os.Getenv("GITHUB_TOKEN")),
	)
	if err != nil {
		t.Fatalf("❌ Failed to initialize Github client (%v)", err)
	}
	ctx := context.Background()

	tags, err := gh.ListAllTags(ctx, vcs.RepositoryAddr{
		Org:  testOrg,
		Name: testRepo,
	})
	if err != nil {
		t.Fatalf("❌ Failed to list GitHub tags (%v)", err)
	}
	for _, tag := range tags {
		if tag.VersionNumber.Equals(testVersion) {
			if !tag.Created.Equal(testDate) {
				t.Fatalf(
					"❌ Found version %s, but the tag date is incorrect: %s (expected: %s)",
					testVersion,
					tag.Created.UTC().String(),
					testDate.UTC().String(),
				)
			}
			t.Logf("✅ Found version %s with tag date %s.", testVersion, tag.Created.String())
			return
		}
	}
	t.Fatalf("❌ Expected version not found (%s)", testVersion)
}

func TestClone(t *testing.T) {
	t.Parallel()

	const testOrg = "integrations"
	const testRepo = "terraform-provider-github"
	const testVersion = "v6.2.3"

	checkoutDir := t.TempDir()

	gh, err := github.New(
		github.WithCheckoutRootDirectory(checkoutDir),
		github.WithLogger(logger.NewTestLogger(t)),
		github.WithToken(os.Getenv("GITHUB_TOKEN")),
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
	t.Parallel()

	var notFound *vcs.RepositoryNotFoundError

	t.Logf("⚙️ Checking if cloning a non-existent repository correctly returns a %T...", notFound)

	const testOrg = "opentofu"
	const testRepo = "nonexistent"
	const testVersion = "v1.6.0"

	checkoutDir := t.TempDir()

	gh, err := github.New(
		github.WithCheckoutRootDirectory(checkoutDir),
		github.WithLogger(logger.NewTestLogger(t)),
		github.WithToken(os.Getenv("GITHUB_TOKEN")),
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

func TestURL(t *testing.T) {
	t.Parallel()

	t.Logf("⚙️ Checking if the GitHub abstraction returns the correct URLs...")
	const testOrg = "opentofu"
	const testRepo = "opentofu"
	const testVersion = "v1.6.0"
	const testFile = "README.md"

	checkoutDir := t.TempDir()
	gh, err := github.New(
		github.WithCheckoutRootDirectory(checkoutDir),
		github.WithLogger(logger.NewTestLogger(t)),
		github.WithToken(os.Getenv("GITHUB_TOKEN")),
	)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	repoURL, err := gh.GetRepositoryBrowseURL(ctx, vcs.RepositoryAddr{Org: testOrg, Name: testRepo})
	if err != nil {
		t.Fatalf("❌ Querying the repo browse URL returned an error (%v)", err)
	}
	if repoURL != "https://github.com/opentofu/opentofu" {
		t.Fatalf("❌ Querying the repo browse URL returned the incorrect URL: %s", repoURL)
	}
	t.Logf("✅ The repo browse URL is correct: %s", repoURL)

	versionURL, err := gh.GetVersionBrowseURL(ctx, vcs.RepositoryAddr{Org: testOrg, Name: testRepo}, testVersion)
	if err != nil {
		t.Fatalf("❌ Querying the repo browse URL returned an error (%v)", err)
	}
	if versionURL != "https://github.com/opentofu/opentofu/tree/v1.6.0" {
		t.Fatalf("❌ Querying the version browse URL returned the incorrect URL: %s", versionURL)
	}
	t.Logf("✅ The version browse URL is correct: %s", versionURL)

	fileURL, err := gh.GetFileViewURL(ctx, vcs.RepositoryAddr{Org: testOrg, Name: testRepo}, testVersion, testFile)
	if err != nil {
		t.Fatalf("❌ Querying the repo view URL returned an error (%v)", err)
	}
	if fileURL != "https://github.com/opentofu/opentofu/blob/v1.6.0/README.md" {
		t.Fatalf("❌ Querying the file view URL returned the incorrect URL: %s", fileURL)
	}
	t.Logf("✅ The file view URL is correct: %s", fileURL)
}
