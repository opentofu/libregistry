// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpg_key_verifier

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func (gv gpgKeyVerifier) ValidateSignature(data string, signature string) error {
	plainMessage := crypto.NewPlainMessageFromString(data)
	pgpSignature, err := crypto.NewPGPSignatureFromArmored(signature)
	if err != nil {
		return fmt.Errorf("failed to create pgp signature: %w", err)
	}

	if err := gv.keyring.VerifyDetached(plainMessage, pgpSignature, crypto.GetUnixTime()); err != nil {
		return fmt.Errorf("failed to verify the detached signature: %w", err)
	}

	return nil
}
