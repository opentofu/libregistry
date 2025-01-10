package key_verification

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func parseKeyRing(pubKeyObj *crypto.Key) (*crypto.KeyRing, error) {
	signingKeyRing, err := crypto.NewKeyRing(pubKeyObj)
	if err != nil {
		return nil, fmt.Errorf("could not build keyring: %w", err)
	}

	return signingKeyRing, nil
}

func validateDetachedSignature(keyring *crypto.KeyRing, data []byte, signature []byte) error {
	plainMessage := crypto.NewPlainMessage(data)
	pgpSignature := crypto.NewPGPSignature(signature)

	if err := keyring.VerifyDetached(plainMessage, pgpSignature, crypto.GetUnixTime()); err != nil {
		return err
	}

	return nil
}
