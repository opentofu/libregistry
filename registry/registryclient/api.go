// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package registryclient

import (
	"context"

	"github.com/opentofu/libregistry/registry/providerregistryprotocol"
)

type Client interface {
	ServiceDiscovery(ctx context.Context) (DiscoveredClient, error)
}

type DiscoveredClient interface {
	providerregistryprotocol.Client
}
