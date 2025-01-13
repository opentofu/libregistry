package gpg_key_verifier

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

type GPGKeyVerifier interface {
	// ValidateSignature validates the signature of the data using the given signature.
	ValidateSignature(data []byte, signature []byte) error
}

type gpgKeyVerifier struct {
	key     *crypto.Key
	keyring *crypto.KeyRing
}

func New(keyData []byte) (GPGKeyVerifier, error) {
	asciiArmor := string(keyData)

	key, err := crypto.NewKeyFromArmored(asciiArmor)
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
