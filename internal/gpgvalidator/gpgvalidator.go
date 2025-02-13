package gpgvalidator

import (
	"context"
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// GPGValidator v
type GPGValidator interface {
	// ValidateSignature validates if the signature was used to sign the data.
	// The keyring used is initialized on New.
	ValidateSignature(ctx context.Context, data []byte, signature []byte) error
}

type gpgValidator struct {
	keyring *crypto.KeyRing
	config  Config
}

func New(key *crypto.Key, options ...Opt) (GPGValidator, error) {
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

	return &gpgValidator{
		keyring: signingKeyRing,
		config:  config,
	}, nil
}
