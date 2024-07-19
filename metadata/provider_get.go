// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) GetProvider(ctx context.Context, providerAddr provider.Addr) (provider.Metadata, error) {
	path := r.getProviderPath(providerAddr)
	fileContents, err := r.storageAPI.GetFile(ctx, path)
	if err != nil {
		var notFound *storage.ErrFileNotFound
		if errors.As(err, &notFound) {
			return provider.Metadata{}, &ProviderNotFoundError{
				ProviderAddr: providerAddr,
				Cause:        err,
			}
		}
		return provider.Metadata{}, fmt.Errorf("failed to read provider file %s (%w)", path, err)
	}
	var mod provider.Metadata
	if err := json.Unmarshal(fileContents, &mod); err != nil {
		return provider.Metadata{}, fmt.Errorf("failed to parse provider metadata file %s (%w)", path, err)
	}
	return mod, nil
}

func (r registryDataAPI) GetProviderCanonicalAddr(_ context.Context, providerAddr provider.Addr) (provider.Addr, error) {
	// TODO implement alias handling
	return providerAddr.Normalize(), nil
}
