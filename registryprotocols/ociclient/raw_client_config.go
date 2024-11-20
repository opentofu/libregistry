// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"github.com/opentofu/libregistry/logger"
	"net/http"
)

type RawOCIClientConfig struct {
	// Credentials contains the per-registry credentials. The client MAY modify this as it obtains a bearer token
	// in exchange for a username-password combination or for a specific OCI name.
	Credentials ScopedCredentials `json:"credentials"`
	// logger contains the logger to use for low-level operations.
	logger logger.Logger
	// HTTPClient contains the HTTP client to use for requests.
	HTTPClient *http.Client
}

func (r *RawOCIClientConfig) ApplyDefaultsAndValidate() error {
	if r.Credentials == nil {
		r.Credentials = map[OCIScopeString]*ClientCredentials{}
	}
	if r.logger == nil {
		r.logger = logger.NewNoopLogger()
	}
	if r.HTTPClient == nil {
		r.HTTPClient = http.DefaultClient
	}
	return nil
}

type RawOCIClientOpt func(c *RawOCIClientConfig) error

func WithRawCredentials(credentials ScopedCredentials) RawOCIClientOpt {
	return func(c *RawOCIClientConfig) error {
		c.Credentials = credentials
		return nil
	}
}

func WithRawHTTPClient(client *http.Client) RawOCIClientOpt {
	return func(c *RawOCIClientConfig) error {
		c.HTTPClient = client
		return nil
	}
}

func WithRawLogger(logger logger.Logger) RawOCIClientOpt {
	return func(c *RawOCIClientConfig) error {
		c.logger = logger
		return nil
	}
}
