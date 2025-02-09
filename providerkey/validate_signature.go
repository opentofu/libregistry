// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"context"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/types/provider"
)

func (pk *providerKey) ValidateSignature(ctx context.Context, pAddr provider.Addr, signature, data []byte) error {
	pk.config.Logger.Info(ctx, "Validating signature with key %s for provider %s...", pk.key.GetHexKeyID(), pAddr.Name)
	plainMessage := crypto.NewPlainMessage(data)
	pgpSignature := crypto.NewPGPSignature(signature)

	if err := pk.config.KeyRing.VerifyDetached(plainMessage, pgpSignature, crypto.GetUnixTime()); err != nil {
		return fmt.Errorf("failed to verify the detached signature: %w", err)
	}

	return nil
}
