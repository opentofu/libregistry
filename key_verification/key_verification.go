package key_verification

import (
	"context"
	"net/http"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage"
)

// KeyVerification describes the functions for verifying if a key was used to sign a list of providers.
type KeyVerification interface {
	// VerifyKey verifies if a key was used to sign a list of providers on the namespace.
	VerifyKey(ctx context.Context, key *crypto.Key, namespace string) error
	// DownloadFile gets an url and return the file contents
	DownloadFile(ctx context.Context, url string) ([]byte, error)
}

// New creates a new instance of the registry API with the given GitHub client and data API instance.
func New(httpClient http.Client, storageAPI storage.API) (KeyVerification, error) {
	dataAPI, err := metadata.New(storageAPI)
	if err != nil {
		return nil, err
	}

	return &keyVerification{
		httpClient: httpClient,
		storageAPI: storageAPI,
		dataAPI:    dataAPI,
	}, nil
}

type keyVerification struct {
	httpClient http.Client
	storageAPI storage.API
	dataAPI    metadata.API
}
