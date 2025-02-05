// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key_verifier

import (
	"context"
	"fmt"

	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

// ProviderKeyVerifier describes the functions for verifying if a key was used to sign a list of providers.
type ProviderKeyVerifier interface {
	// VerifyProvider verifies if the key was used to sign a provider addr. It returns a list of the valid versions signed by this key.
	VerifyProvider(ctx context.Context, provider provider.Addr) ([]*provider.Version, error)
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
		Config:      config,
		gpgVerifier: gpgVerifier,
		dataAPI:     dataAPI,
	}, nil
}

type providerKeyVerifier struct {
	Config      Config
	gpgVerifier gpg_key_verifier.GPGKeyVerifier
	dataAPI     metadata.API
}
