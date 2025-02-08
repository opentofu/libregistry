// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/logger"
)

// Config holds the configuration for GitHub.
type Config struct {
	// Logger holds the logger to write any logs to.
	Logger logger.Logger
	// HTTPClient holds the HTTP client to use for API requests. Note that this only affects API and RSS feed requests,
	// but not git clone commands as those are done using the command line.
	HTTPClient *http.Client
	// Number of versions that are going to be checked if they were signed
	VersionsToCheck uint8
	// Number of max parallelism used when checking the signatures for the versions
	MaxParallelism uint8
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

	if c.MaxParallelism == 0 {
		c.MaxParallelism = 10
	}

	if c.VersionsToCheck == 0 {
		c.VersionsToCheck = 10
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

// WithVersionsToCheck is a functional option to set the number of versions to check for a provider.
func WithNumVersionsToCheck(versionsToCheck uint8) Opt {
	return func(config *Config) error {
		config.VersionsToCheck = versionsToCheck
		return nil
	}
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

// WithMaxParallelism is a functional option to set the function used to check the provider version
func WithMaxParallelism(maxParallelism uint8) Opt {
	return func(config *Config) error {
		config.MaxParallelism = maxParallelism
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
