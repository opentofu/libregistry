// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package servicediscovery

import (
	"context"
	"github.com/opentofu/libregistry/branding"
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
	cfg.applyDefaultsAndValidate()
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

func WithHostname(hostname string) ClientOpt {
	return func(cfg *config) error {
		cfg.hostname = hostname
		return nil
	}
}

func WithScheme(scheme string) ClientOpt {
	return func(cfg *config) error {
		cfg.scheme = scheme
		return nil
	}
}

type config struct {
	scheme   string
	hostname string
	client   *http.Client
}

func (c *config) applyDefaultsAndValidate() error {
	if c.client == nil {
		c.client = http.DefaultClient
	}

	if c.hostname == "" {
		c.hostname = branding.DefaultRegistry
	}

	if c.scheme == "" {
		c.scheme = branding.DefaultRegistryScheme
	}
	return nil
}

type client struct {
	cfg config
}

func (c client) ServiceDiscovery(ctx context.Context, _ Request) (Response, error) {
	return tofuhttp.GetRequest[Response](ctx, c.cfg.client, c.cfg.scheme+"://"+c.cfg.hostname, ".well-known/terraform.json")
}
