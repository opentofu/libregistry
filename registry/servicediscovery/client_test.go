// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package servicediscovery_test

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/registry/servicediscovery"
)

func TestSD(t *testing.T) {
	cli, err := servicediscovery.NewClient()
	if err != nil {
		t.Fatalf("Failed to initialize client (%v)", err)
	}
	resp, err := cli.ServiceDiscovery(context.Background(), servicediscovery.Request{})
	if err != nil {
		t.Fatalf("Failed to perform service discovery (%v)", err)
	}
	if resp.ProvidersV1 != "/v1/providers/" {
		t.Fatalf("Unexpected providers endpoint: %s", resp.ProvidersV1)
	}
}
