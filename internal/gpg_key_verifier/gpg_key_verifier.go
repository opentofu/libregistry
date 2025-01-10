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
	keyring *crypto.KeyRing
}

func New(key *crypto.Key) (GPGKeyVerifier, error) {
	signingKeyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		return nil, fmt.Errorf("could not build GPG verifier: %w", err)
	}

	return &gpgKeyVerifier{
		keyring: signingKeyRing,
	}, nil
}
