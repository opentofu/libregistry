package key_verification

import (
	"context"
	"net/http"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage"
)

// KeyVerification describes the functions for verifying the authenticity of a provider's key.
type KeyVerification interface {
	// AddModule adds a module based on a VCS repository. The VCS repository name must follow the naming convention
	// of the VCS implementation passed to the registry API on initialization.
	VerifyKey(ctx context.Context, keyPath string, namespace string) error
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
	dataAPI    metadata.ProviderDataAPI
}
