// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"context"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

// ProviderKey describes the functions for verifying if a key was used to sign a list of providers.
type ProviderKey interface {
	// VerifyProvider verifies if the key was used to sign a provider addr. It returns a list of the valid versions signed by this key.
	VerifyProvider(ctx context.Context, provider provider.Addr) ([]provider.Version, error)
	// ValidateSignature validates if the signature was used to sign the data. The keyring used is initialized on New.
	ValidateSignature(signature, data []byte) error
}

// New creates a new instance of the provider key verification package with the given f keyData (GPG ASCII-Armored PEM) and the metadata API. There are a few optional fields that can be used modify the behavior of the package.
func New(keyData string, dataAPI metadata.API, options ...Opt) (ProviderKey, error) {
	key, err := crypto.NewKeyFromArmored(keyData)
	if err != nil {
		return nil, fmt.Errorf("could not parse armored key: %w", err)
	}

	config := Config{}
	for _, opt := range options {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	err = config.ApplyDefaults(key)
	if err != nil {
		return nil, fmt.Errorf("failed to apply defaults: %w", err)
	}

	return &providerKey{
		config:  config,
		dataAPI: dataAPI,
	}, nil
}

type providerKey struct {
	config  Config
	dataAPI metadata.API
}
