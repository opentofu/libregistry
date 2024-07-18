package metadata

import (
	"github.com/opentofu/libregistry/metadata/storage"
)

// API providers the methods for accessing the stored registry data.
type API interface {
	ModuleDataAPI
	ProviderDataAPI
}

// New creates a new API.
func New(storageAPI storage.API) (API, error) {
	return &registryDataAPI{
		storageAPI: storageAPI,
	}, nil
}

type registryDataAPI struct {
	storageAPI storage.API
}
