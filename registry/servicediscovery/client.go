// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package servicediscovery

import (
	"context"
	"net/http"

	tofuhttp "github.com/opentofu/libregistry/registry/internal/http"
)

type Client interface {
	ServiceDiscovery(ctx context.Context, request Request) (Response, error)
}

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

func (c client) ServiceDiscovery(ctx context.Context, _ Request) (Response, error) {
	return tofuhttp.GetRequest[Response](ctx, c.cfg.client, c.cfg.endpoint, ".well-known/terraform.json")
}
