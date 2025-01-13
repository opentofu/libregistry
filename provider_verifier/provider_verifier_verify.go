package provider_verifier

import (
	"context"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func (kv keyVerification) VerifyKey(ctx context.Context, key *crypto.Key, providerAddr provider.Addr) error {
	gpgVerifier, err := gpg_key_verifier.New(key)
	if err != nil {
		return fmt.Errorf("failed to verify key %s for provider %s (cannot construct GPG key verifier: %w)", key.GetHexKeyID(), providerAddr, err)
	}

	provider, err := kv.dataAPI.GetProvider(ctx, providerAddr, false)
	if err != nil {
		return fmt.Errorf("failed to verify key %s for provider %s (%w)", key.GetHexKeyID(), providerAddr, err)
	}

	for _, version := range provider.Versions {
		shaSumContents, err := kv.DownloadFile(ctx, version.SHASumsURL)
		if err != nil {
			return fmt.Errorf("failed to verify key %s for provider %s (%w)", key.GetHexKeyID(), providerAddr, err)
		}

		shaSumSigContents, err := kv.DownloadFile(ctx, version.SHASumsSignatureURL)
		if err != nil {
			return fmt.Errorf("failed to verify key %s for provider %s (%w)", key.GetHexKeyID(), providerAddr, err)
		}

		if err := gpgVerifier.ValidateSignature(shaSumContents, shaSumSigContents); err != nil {
			return fmt.Errorf("failed to verify key %s for provider %s (%w)", key.GetHexKeyID(), providerAddr, err)
		}
	}

	return nil
}
