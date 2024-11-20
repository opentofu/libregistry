// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
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
	if err != nil {
		t.Fatalf("%v", err)
	}
	for _, warning := range warnings {
		t.Logf("OCI registry warning: %s", warning)
	}
	t.Logf("Found the following references: %v", references)
	for _, ref := range references {
		if ref == "1.6.0" {
			t.Logf("Version 1.6.0 is present in the version list.")
			return
		}
	}
	t.Fatalf("Version 1.6.0 was not found in the version list.")
}
