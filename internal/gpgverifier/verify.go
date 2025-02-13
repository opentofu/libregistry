// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpgverifier

import (
	"context"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func (gv *gpgVerifier) VerifySignature(ctx context.Context, signature, data []byte) error {
	plainMessage := crypto.NewPlainMessage(data)
	pgpSignature := crypto.NewPGPSignature(signature)

	if err := gv.keyring.VerifyDetached(plainMessage, pgpSignature, crypto.GetUnixTime()); err != nil {
		return fmt.Errorf("failed to verify the detached signature: %w", err)
	}

	return nil
}
