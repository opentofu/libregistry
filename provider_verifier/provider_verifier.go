package provider_verifier

import (
	"context"
	"net/http"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

// KeyVerification describes the functions for verifying if a key was used to sign a list of providers.
type KeyVerification interface {
	// VerifyKey verifies if a key was used to sign a list of providers on the namespace.
	VerifyKey(ctx context.Context, keyData []byte, provider provider.Addr) error
	// downloadFile gets an url and return the file contents
	downloadFile(ctx context.Context, url string) ([]byte, error)
}

// New creates a new instance of the key verification package with the given http client and a storage instance.
func New(httpClient http.Client, dataAPI metadata.API) (KeyVerification, error) {
	return &keyVerification{
		httpClient: httpClient,
		dataAPI:    dataAPI,
	}, nil
}

type keyVerification struct {
	httpClient http.Client
	dataAPI    metadata.API
}
