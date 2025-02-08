package provider_key

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func (pk *providerKey) ValidateSignature(signature, data []byte) error {
	plainMessage := crypto.NewPlainMessage(data)
	pgpSignature := crypto.NewPGPSignature(signature)

	if err := pk.config.KeyRing.VerifyDetached(plainMessage, pgpSignature, crypto.GetUnixTime()); err != nil {
		return fmt.Errorf("failed to verify the detached signature: %w", err)
	}

	return nil
}
