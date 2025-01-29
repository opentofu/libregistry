package provider_key_verifier

import (
	"context"
	"fmt"

	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func (pkv providerKeyVerifier) VerifyKey(ctx context.Context, keyData []byte, providerAddr provider.Addr) error {
	gpgVerifier, err := gpg_key_verifier.New(keyData)
	if err != nil {
		return fmt.Errorf("failed to verify key for provider %s (cannot construct GPG key verifier: %w)", providerAddr, err)
	}

	provider, err := pkv.dataAPI.GetProvider(ctx, providerAddr, false)
	if err != nil {
		return fmt.Errorf("failed to get provider %s (%w)", providerAddr, err)
	}

	toCheck := min(uint8(len(provider.Versions)), pkv.versionsToCheck)

	for _, version := range provider.Versions[:toCheck] {
		shaSumContents, err := pkv.downloadFile(ctx, version.SHASumsURL)
		if err != nil {
			return fmt.Errorf("failed to download SHASums URL for provider %s (%w)", providerAddr, err)
		}

		shaSumSigContents, err := pkv.downloadFile(ctx, version.SHASumsSignatureURL)
		if err != nil {
			return fmt.Errorf("failed to download SHASums signature URL for provider %s (%w)", providerAddr, err)
		}

		if err := gpgVerifier.ValidateSignature(shaSumContents, shaSumSigContents); err != nil {
			return fmt.Errorf("failed to validate signature for provider %s (%w)", providerAddr, err)
		}
	}

	return nil
}
