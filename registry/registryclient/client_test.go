// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package registryclient_test

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/registry/providerregistry"
	"github.com/opentofu/libregistry/registry/registryclient"
)

func TestClient(t *testing.T) {
	cli, err := registryclient.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client (%v)", err)
	}
	ctx := context.Background()

	registrycli, err := cli.ServiceDiscovery(ctx)
	if err != nil {
		t.Fatalf("Service discovery failed (%v)", err)
	}

	awsResponse, err := registrycli.ListAvailableProviderVersions(ctx, providerregistry.ListAvailableProviderVersionsRequest{
		Namespace: "opentofu",
		Type:      "aws",
	})
	if err != nil {
		t.Fatalf("Failed to fetch the provider versions (%v)", err)
	}
	if len(awsResponse.Versions) == 0 {
		t.Fatalf("Response contained no versions")
	}
	found := false
	for _, ver := range awsResponse.Versions {
		if ver.Version == "5.65.0" {
			t.Logf("Found version 5.65.0.")
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("No version 5.65.0 found.")
	}

	version, err := registrycli.GetProviderVersion(ctx, providerregistry.GetProviderVersionRequest{
		Namespace: "opentofu",
		Type:      "aws",
		Version:   "5.65.0",
		OS:        "darwin",
		Arch:      "amd64",
	})
	if err != nil {
		t.Fatalf("Failed to fetch provider version (%v)", err)
	}
	t.Logf("Download URL is: %s", version.DownloadURL)
}
