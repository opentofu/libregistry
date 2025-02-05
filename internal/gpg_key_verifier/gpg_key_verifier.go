// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpg_key_verifier

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// GPGKeyVerifier provides functions for validating GPG Signatures. It's mostly used for the provider keys used to sign SHASUM urls, but it can be used for any place that GPG keys were used.
type GPGKeyVerifier interface {
	// ValidateSignature validates the signature of the data using the given signature.
	ValidateSignature(data []byte, signature []byte) error
}

type gpgKeyVerifier struct {
	key     *crypto.Key
	keyring *crypto.KeyRing
}

// New creates a GPG Key Verifier by passing ASCII-Armored PEM data in the keyData attribute.
func New(keyData string) (GPGKeyVerifier, error) {
	key, err := crypto.NewKeyFromArmored(keyData)
	if err != nil {
		return nil, fmt.Errorf("could not parse armored key: %w", err)
	}

	signingKeyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		return nil, fmt.Errorf("could not build keyring for key %s: %w", key.GetHexKeyID(), err)
	}

	return &gpgKeyVerifier{
		key:     key,
		keyring: signingKeyRing,
	}, nil
}
