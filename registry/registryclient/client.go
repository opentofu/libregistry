// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package registryclient

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/opentofu/libregistry/registry/providerregistryprotocol"
	"github.com/opentofu/libregistry/registry/servicediscovery"
)

type ClientOpt func(cfg *config) error

func NewClient(opts ...ClientOpt) (Client, error) {
	cfg := config{}
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}
	cfg.applyDefaults()
	return &client{
		cfg,
	}, nil
}

func WithHTTPClient(client *http.Client) ClientOpt {
	return func(cfg *config) error {
		cfg.client = client
		return nil
	}
}

func WithEndpoint(endpoint string) ClientOpt {
	return func(cfg *config) error {
		cfg.endpoint = endpoint
		return nil
	}
}

type config struct {
	endpoint string
	client   *http.Client
}

func (c *config) applyDefaults() {
	if c.client == nil {
		c.client = http.DefaultClient
	}
	if c.endpoint == "" {
		c.endpoint = "https://registry.opentofu.org"
	}
}

type client struct {
	cfg config
}

func (c client) ServiceDiscovery(ctx context.Context) (DiscoveredClient, error) {
	sd, err := servicediscovery.NewClient(
		servicediscovery.WithEndpoint(c.cfg.endpoint),
		servicediscovery.WithHTTPClient(c.cfg.client),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize service discovery client (%w)", err)
	}
	resp, err := sd.ServiceDiscovery(ctx, servicediscovery.Request{})
	if err != nil {
		return nil, fmt.Errorf("failed to perform service discovery (%w)", err)
	}

	if resp.ProvidersV1 == "" {
		return nil, fmt.Errorf("no providers endpoint found")
	}

	return providerregistryprotocol.NewClient(
		providerregistryprotocol.WithHTTPClient(c.cfg.client),
		providerregistryprotocol.WithProvidersEndpoint(strings.TrimSuffix(c.cfg.endpoint, "/")+resp.ProvidersV1),
	)
}
