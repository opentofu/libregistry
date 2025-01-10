package key_verification

import (
	"context"
	"fmt"
	"os"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (kv keyVerification) VerifyKey(ctx context.Context, keyPath string, namespace string) error {
	namespace = provider.NormalizeNamespace(namespace)
	keyFile, err := kv.storageAPI.GetFile(ctx, storage.Path(keyPath))

	if err != nil {
		return fmt.Errorf("failed to load the key %s (%w)", keyPath, err)
	}

	signingKeyRing, err := parseSigningKeyRing(string(keyFile))

	providers, err := kv.dataAPI.ListProvidersByNamespace(ctx, namespace, false)

	for _, providerAddr := range providers {
		provider, err := kv.dataAPI.GetProvider(ctx, providerAddr, false)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(err.Error()))
		}

		for _, version := range provider.Versions {
			shaSumContents, err := kv.DownloadFile(ctx, version.SHASumsURL)
			if err != nil {
				return err
			}

			shaSumSigContents, err := kv.DownloadFile(ctx, version.SHASumsSignatureURL)
			if err != nil {
				return err
			}

			if err := validateDetachedSignature(signingKeyRing, shaSumContents, shaSumSigContents); err != nil {
				return err
			}
		}
	}
	return nil
}
