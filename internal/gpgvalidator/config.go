// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpgvalidator

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/logger"
)

// Config holds the configuration for ProviderKey.
type Config struct {
	// Logger holds the logger to write any logs to.
	Logger logger.Logger
	// HTTPClient holds the HTTP client to use for API requests. Note that this only affects API and RSS feed requests,
	// but not git clone commands as those are done using the command line.
	HTTPClient *http.Client
	// Keyring is used to test the PGP key
	KeyRing *crypto.KeyRing
}

// Opt is a function that modifies the config.
type Opt func(config *Config) error

// ApplyDefaults adds the default values if none are present.
func (c *Config) ApplyDefaults(key *crypto.Key) error {
	if c.Logger == nil {
		c.Logger = logger.NewNoopLogger()
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
		transport := http.DefaultTransport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
		c.HTTPClient.Transport = transport
	}

	if c.KeyRing == nil {
		keyring, err := crypto.NewKeyRing(key)
		if err != nil {
			return fmt.Errorf("could not build keyring for key %s: %w", key.GetHexKeyID(), err)
		}
		c.KeyRing = keyring
	}

	return nil
}

// WithLogger is a functional option to set the logger
func WithLogger(logger logger.Logger) Opt {
	return func(config *Config) error {
		config.Logger = logger
		return nil
	}
}

// WithHTTPClient is a functional option to set the http Client
func WithHTTPClient(httpClient *http.Client) Opt {
	return func(config *Config) error {
		config.HTTPClient = httpClient
		return nil
	}
}

// WithKeyring allows to define a PGP KeyRing
func WithKeyring(keyring *crypto.KeyRing) Opt {
	return func(config *Config) error {
		config.KeyRing = keyring
		return nil
	}
}
