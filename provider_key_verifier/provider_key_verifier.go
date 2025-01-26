package provider_key_verifier

import (
	"context"
	"net/http"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

// ProviderKeyVerifier describes the functions for verifying if a key was used to sign a list of providers.
type ProviderKeyVerifier interface {
	// VerifyKey verifies if a keyData (GPG ASCII-Armored PEM) was used to sign a provider addr.
	VerifyKey(ctx context.Context, keyData []byte, provider provider.Addr) error
}

// New creates a new instance of the key verification package with the given http client and a metadata instance.
func New(httpClient http.Client, dataAPI metadata.API, opts ...Option) (ProviderKeyVerifier, error) {
	providerKeyVerifier := &providerKeyVerifier{
		httpClient:      httpClient,
		dataAPI:         dataAPI,
		versionsToCheck: 10,
	}

	for _, opt := range opts {
		opt(providerKeyVerifier)
	}

	return providerKeyVerifier, nil
}

// Option is used for providing options to New without changing the signature of New.
type Option func(c *providerKeyVerifier)

// WithVersionsToCheck is a functional option to set the number of versions to check for a provider.
func WithVersionsToCheck(versionsToCheck uint8) Option {
	return func(c *providerKeyVerifier) {
		c.versionsToCheck = versionsToCheck
	}
}

type providerKeyVerifier struct {
	httpClient      http.Client
	dataAPI         metadata.API
	versionsToCheck uint8
}
