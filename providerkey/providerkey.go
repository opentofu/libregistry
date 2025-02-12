// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

// Package providerkey is used to verify if GPG key were used to sign
// OpenTofu providers, using its version and shaSumsURL and shaSumsSignatureURL.
package providerkey

import (
	"context"
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

// ProviderKey describes the functions for verifying if a key was used to sign a list of providers.
type ProviderKey interface {
	// VerifyProvider verifies if the key was used to sign a provider addr.
	// It returns a list of the valid versions signed by this key.
	VerifyProvider(ctx context.Context, provider provider.Addr) ([]provider.Version, error)
	// ValidateSignature validates if the signature was used to sign the data.
	// The keyring used is initialized on New.
	ValidateSignature(ctx context.Context, provider provider.Addr, signature, data []byte) error
}

// New creates a new instance of the provider key verification package with
// the given f keyData (GPG ASCII-Armored PEM) and the metadata API.
// There are a few optional fields that can be used modify the behavior of the package.
func New(keyData string, dataAPI metadata.API, options ...Opt) (ProviderKey, error) {
	key, err := crypto.NewKeyFromArmored(keyData)
	if err != nil {
		return nil, fmt.Errorf("could not parse armored key: %w", err)
	}

	config := Config{}
	var errs error
	for _, opt := range options {
		if err := opt(&config); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, fmt.Errorf("failed to apply config options: %w", err)
	}

	err = config.ApplyDefaults(key)
	if err != nil {
		return nil, fmt.Errorf("failed to apply defaults: %w", err)
	}

	return &providerKey{
		key:     key,
		config:  config,
		dataAPI: dataAPI,
	}, nil
}

type providerKey struct {
	config  Config
	key     *crypto.Key
	dataAPI metadata.API
}
