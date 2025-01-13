package provider_verifier

import (
	"context"
	"fmt"

	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func (kv keyVerification) VerifyKey(ctx context.Context, keyData []byte, providerAddr provider.Addr, versionsToCheck uint16) error {
	if versionsToCheck == 0 {
		versionsToCheck = 10
	}

	gpgVerifier, err := gpg_key_verifier.New(keyData)
	if err != nil {
		return fmt.Errorf("failed to verify key for provider %s (cannot construct GPG key verifier: %w)", providerAddr, err)
	}

	provider, err := kv.dataAPI.GetProvider(ctx, providerAddr, false)
	if err != nil {
		return fmt.Errorf("failed to get provider %s (%w)", providerAddr, err)
	}

	for _, version := range provider.Versions[:versionsToCheck] {
		shaSumContents, err := downloadFile(ctx, kv.httpClient, version.SHASumsURL)
		if err != nil {
			return fmt.Errorf("failed to download SHASums URL for provider %s (%w)", providerAddr, err)
		}

		shaSumSigContents, err := downloadFile(ctx, kv.httpClient, version.SHASumsSignatureURL)
		if err != nil {
			return fmt.Errorf("failed to download SHASums signature URL for provider %s (%w)", providerAddr, err)
		}

		if err := gpgVerifier.ValidateSignature(shaSumContents, shaSumSigContents); err != nil {
			return fmt.Errorf("failed to validate signature for provider %s (%w)", providerAddr, err)
		}
	}

	return nil
}
