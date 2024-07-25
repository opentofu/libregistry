package github_test

import (
	"context"
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
