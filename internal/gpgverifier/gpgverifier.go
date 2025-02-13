// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpgverifier

import (
	"context"
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// GPGVerifier validates if the data was signed by the binary signature.
// data comes from the *_SHA256SUMS file and signature from _SHA256SUMS.sig file from the provider.
type GPGVerifier interface {
	// VerifySignature validates if the signature was used to sign the data.
	// The keyring used is initialized on New.
	VerifySignature(ctx context.Context, data []byte, signature []byte) error
}

type gpgVerifier struct {
	keyring *crypto.KeyRing
	config  Config
}

func New(key *crypto.Key, options ...Opt) (GPGVerifier, error) {
	signingKeyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		return nil, fmt.Errorf("could not build GPG verifier: %w", err)
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

	return &gpgVerifier{
		keyring: signingKeyRing,
		config:  config,
	}, nil
}
