// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package registryclient

import (
	"context"
	"fmt"
	"github.com/opentofu/libregistry/branding"
	"net/http"
	"strings"

	"github.com/opentofu/libregistry/registry/providerregistry"
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
	if err := cfg.applyDefaultsAndValidate(); err != nil {
		return nil, err
	}
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

func WithScheme(scheme string) ClientOpt {
	return func(cfg *config) error {
		if scheme != "http" && scheme != "https" {
			return fmt.Errorf("invalid scheme: %s (must be https or http)", scheme)
		}
		cfg.scheme = scheme
		return nil
	}
}

func WithHostname(endpoint string) ClientOpt {
	return func(cfg *config) error {
		cfg.hostname = endpoint
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
	if c.scheme == "" {
		c.scheme = "https"
	}
	if c.hostname == "" {
		c.hostname = branding.DefaultRegistry
	}
	return nil
}

type client struct {
	cfg config
}

func (c client) ServiceDiscovery(ctx context.Context) (DiscoveredClient, error) {
	sd, err := servicediscovery.NewClient(
		servicediscovery.WithScheme(c.cfg.scheme),
		servicediscovery.WithHostname(c.cfg.hostname),
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

	var providersEndpoint string
	if strings.HasPrefix(resp.ProvidersV1, "/") {
		providersEndpoint = c.cfg.scheme + "://" + c.cfg.hostname + resp.ProvidersV1
	} else {
		providersEndpoint = resp.ProvidersV1
	}

	return providerregistry.NewClient(
		providerregistry.WithHTTPClient(c.cfg.client),
		providerregistry.WithProvidersEndpoint(providersEndpoint),
	)
}
