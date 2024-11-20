// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient_test

import (
	"context"
	"github.com/opentofu/libregistry/logger"
	"github.com/opentofu/libregistry/registryprotocols/ociclient"
	"testing"
)

func TestListReferences(t *testing.T) {
	client, err := ociclient.New(ociclient.WithLogger(logger.NewTestLogger(t)))
	if err != nil {
		t.Fatalf("%v", err)
	}
	references, warnings, err := client.ListReferences(context.Background(), ociclient.OCIAddr{
		Registry: "ghcr.io",
		Name:     "opentofu/opentofu",
	})
	for _, warning := range warnings {
		t.Logf("OCI registry warning: %s", warning)
	}
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("Found the following references: %v", references)
	for _, ref := range references {
		if ref == "latest" {
			t.Logf("Version latest is present in the version list.")
			return
		}
	}
	t.Fatalf("Version latest was not found in the version list.")
}

func TestPull(t *testing.T) {
	client, err := ociclient.New(ociclient.WithLogger(logger.NewTestLogger(t)))
	if err != nil {
		t.Fatalf("%v", err)
	}
	image, warnings, err := client.PullImage(
		context.Background(),
		ociclient.OCIAddrWithReference{
			OCIAddr: ociclient.OCIAddr{
				Registry: "ghcr.io",
				Name:     "opentofu/opentofu",
			},
			Reference: "latest",
		},
		ociclient.WithGOOS("linux"),
		ociclient.WithGOARCH("amd64"),
	)
	for _, warning := range warnings {
		t.Logf("OCI registry warning: %s", warning)
	}
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer func() {
		if err := image.Close(); err != nil {
			t.Fatalf("%v", err)
		}
	}()
	found := false
	for {
		ok, err := image.Next()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if !ok {
			break
		}
		t.Logf("Found file: %s", image.Filename())
		if image.Filename() == "usr/local/bin/tofu" {
			found = true
		}
	}
	if !found {
		t.Fatalf("no tofu found in downloaded image")
	}
}
