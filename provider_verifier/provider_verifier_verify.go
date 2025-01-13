package provider_verifier

import (
	"context"
	"fmt"

	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func (kv keyVerification) VerifyKey(ctx context.Context, keyData []byte, providerAddr provider.Addr) error {
	gpgVerifier, err := gpg_key_verifier.New(keyData)
	if err != nil {
		return fmt.Errorf("failed to verify key %s for provider %s (cannot construct GPG key verifier: %w)", gpgVerifier.GetHexKeyID(), providerAddr, err)
	}

	provider, err := kv.dataAPI.GetProvider(ctx, providerAddr, false)
	if err != nil {
		return fmt.Errorf("failed to verify key %s for provider %s (%w)", gpgVerifier.GetHexKeyID(), providerAddr, err)
	}

	for _, version := range provider.Versions {
		shaSumContents, err := downloadFile(ctx, kv.httpClient, version.SHASumsURL)
		if err != nil {
			return fmt.Errorf("failed to verify key %s for provider %s (%w)", gpgVerifier.GetHexKeyID(), providerAddr, err)
		}

		shaSumSigContents, err := downloadFile(ctx, kv.httpClient, version.SHASumsSignatureURL)
		if err != nil {
			return fmt.Errorf("failed to verify key %s for provider %s (%w)", gpgVerifier.GetHexKeyID(), providerAddr, err)
		}

		if err := gpgVerifier.ValidateSignature(shaSumContents, shaSumSigContents); err != nil {
			return fmt.Errorf("failed to verify key %s for provider %s (%w)", gpgVerifier.GetHexKeyID(), providerAddr, err)
		}
	}

	return nil
}
