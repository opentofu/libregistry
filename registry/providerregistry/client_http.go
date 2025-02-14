// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerregistry

import (
	"context"
	"github.com/opentofu/libregistry/branding"
	"net/http"
	"net/url"

	tofuhttp "github.com/opentofu/libregistry/registry/internal/http"
)

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

type ClientOpt func(cfg *config) error

func WithHTTPClient(client *http.Client) ClientOpt {
	return func(cfg *config) error {
		cfg.client = client
		return nil
	}
}

func WithProvidersEndpoint(endpoint string) ClientOpt {
	return func(cfg *config) error {
		cfg.providersEndpoint = endpoint
		return nil
	}
}

type config struct {
	providersEndpoint string
	client            *http.Client
}

func (c *config) applyDefaultsAndValidate() error {
	if c.client == nil {
		c.client = http.DefaultClient
	}
	if c.providersEndpoint == "" {
		c.providersEndpoint = branding.DefaultProvidersEndpoint
	}
	return nil
}

type client struct {
	cfg config
}

func (c client) GetProviderVersion(ctx context.Context, request GetProviderVersionRequest) (response GetProviderVersionResponse, err error) {
	urlSuffix := url.PathEscape(string(request.Namespace)) + "/" + url.PathEscape(string(request.Type)) + "/" + url.PathEscape(string(request.Version)) + "/download/" + url.PathEscape(string(request.OS)) + "/" + url.PathEscape(string(request.Arch))
	return tofuhttp.GetRequest[GetProviderVersionResponse](ctx, c.cfg.client, c.cfg.providersEndpoint, urlSuffix)
}

func (c client) ListAvailableProviderVersions(ctx context.Context, request ListAvailableProviderVersionsRequest) (ListAvailableProviderVersionsResponse, error) {
	urlSuffix := url.PathEscape(string(request.Namespace)) + "/" + url.PathEscape(string(request.Type)) + "/versions"
	return tofuhttp.GetRequest[ListAvailableProviderVersionsResponse](ctx, c.cfg.client, c.cfg.providersEndpoint, urlSuffix)
}
