// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerregistry_test

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/registry/providerregistry"
)

func TestHTTPClient(t *testing.T) {
	cli, err := providerregistry.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client (%v)", err)
	}
	ctx := context.Background()

	awsResponse, err := cli.ListAvailableProviderVersions(ctx, providerregistry.ListAvailableProviderVersionsRequest{
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

	version, err := cli.GetProviderVersion(ctx, providerregistry.GetProviderVersionRequest{
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
