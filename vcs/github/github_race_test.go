// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package github_test

import (
	"context"
	"os"
	"testing"

	"github.com/opentofu/libregistry/logger"
	"github.com/opentofu/libregistry/vcs"
	"github.com/opentofu/libregistry/vcs/github"
	"golang.org/x/sync/errgroup"
)

func TestCloneRace(t *testing.T) {
	t.Parallel()

	const testParallelism = 10

	const testOrg = "opentofu"
	const testRepo = "terraform-provider-tfcoremock"
	const testVersion = "v0.3.0"

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
		t.Fatalf("%v", err)
	}
	if err := workingCopy.Close(); err != nil {
		t.Fatalf("%v", err)
	}

	start := make(chan struct{})

	errGroup := &errgroup.Group{}
	for i := 0; i < testParallelism; i++ {
		errGroup.Go(func() error {
			<-start
			workingCopy, err := gh.Checkout(ctx, vcs.RepositoryAddr{
				Org:  testOrg,
				Name: testRepo,
			}, testVersion)
			if err != nil {
				return err
			}
			return workingCopy.Close()
		})
	}
	close(start)
	if err := errGroup.Wait(); err != nil {
		t.Fatalf("%v", err)
	}
}
