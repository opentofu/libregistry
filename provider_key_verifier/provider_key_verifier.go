// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key_verifier

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
	"github.com/opentofu/libregistry/logger"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

// ProviderKeyVerifier describes the functions for verifying if a key was used to sign a list of providers.
type ProviderKeyVerifier interface {
	// VerifyProvider verifies if the key was used to sign a provider addr. It returns a list of the valid versions signed by this key.
	VerifyProvider(ctx context.Context, provider provider.Addr) ([]*provider.Version, error)
}

// Config holds the configuration for GitHub.
type Config struct {
	// Logger holds the logger to write any logs to.
	Logger logger.Logger
	// HTTPClient holds the HTTP client to use for API requests. Note that this only affects API and RSS feed requests,
	// but not git clone commands as those are done using the command line.
	HTTPClient *http.Client
	// Number of versions that are going to be checked if they were signed
	NumVersionsToCheck uint8
	checkFn            CheckFn
}

// New creates a new instance of the provider key verification package with the given f keyData (GPG ASCII-Armored PEM) and the metadata API. There are a few optional fields that can be used modify the behavior of the package.
func New(keyData string, dataAPI metadata.API, options ...Opt) (ProviderKeyVerifier, error) {
	gpgVerifier, err := gpg_key_verifier.New(keyData)
	if err != nil {
		return nil, fmt.Errorf("cannot construct GPG key verifier: %w", err)
	}

	config := Config{}
	for _, opt := range options {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}
	config.ApplyDefaults()

	return &providerKeyVerifier{
		config:      config,
		gpgVerifier: gpgVerifier,
		dataAPI:     dataAPI,
	}, nil
}

type providerKeyVerifier struct {
	gpgVerifier gpg_key_verifier.GPGKeyVerifier
	dataAPI     metadata.API
	config      Config
}
