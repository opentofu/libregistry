package key_verification

import (
	"context"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/internal/gpg_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func (kv keyVerification) VerifyKey(ctx context.Context, key *crypto.Key, namespace string) error {
	namespace = provider.NormalizeNamespace(namespace)

	gpgVerifier, err := gpg_verifier.New(key)
	if err != nil {
		return err
	}

	providers, err := kv.dataAPI.ListProvidersByNamespace(ctx, namespace, false)
	if err != nil {
		return err
	}

	for _, providerAddr := range providers {
		provider, err := kv.dataAPI.GetProvider(ctx, providerAddr, false)
		if err != nil {
			return err
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

			if err := gpgVerifier.ValidateSignature(shaSumContents, shaSumSigContents); err != nil {
				return err
			}
		}
	}
	return nil
}
