// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/opentofu/libregistry/internal/gpg"
	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) CheckProviderNamespaceKey(ctx context.Context, namespace string, keyPath string) error {
	namespace = provider.NormalizeNamespace(namespace)
	keyFile, err := r.storageAPI.GetFile(ctx, storage.Path(keyPath))

	if err != nil {
		return fmt.Errorf("failed to load the key %s (%w)", keyPath, err)
	}

	signingKeyRing, err := gpg.ParseSigningKeyRing(string(keyFile))

	providers, err := r.ListProvidersByNamespace(ctx, namespace, false)

	for _, providerAddr := range providers {
		provider, err := r.GetProvider(ctx, providerAddr, false)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(err.Error()))
		}

		for _, version := range provider.Versions {
			shaSumsTmpPath := storage.Path(fmt.Sprintf("/tmp/%s", path.Base(version.SHASumsURL)))
			shaSumsSigTmpPath := storage.Path(fmt.Sprintf("/tmp/%s", path.Base(version.SHASumsSignatureURL)))

			if err := r.storageAPI.DownloadFile(ctx, version.SHASumsSignatureURL, shaSumsSigTmpPath); err != nil {
				return err
			}

			if err := r.storageAPI.DownloadFile(ctx, version.SHASumsURL, shaSumsTmpPath); err != nil {
				return err
			}

			shaSumContents, err := r.storageAPI.GetFile(ctx, shaSumsTmpPath)
			if err != nil {
				return err
			}

			shaSumSigContents, err := r.storageAPI.GetFile(ctx, shaSumsSigTmpPath)
			if err != nil {
				return err
			}

			if err := gpg.ValidateDetachedSignature(signingKeyRing, shaSumContents, shaSumSigContents); err != nil {
				return err
			}
		}
	}
	return nil
}
