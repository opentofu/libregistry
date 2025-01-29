package provider_key_verifier

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func (pkv providerKeyVerifier) VerifyKey(ctx context.Context, keyData []byte, providerAddr provider.Addr) ([]string, error) {
	gpgVerifier, err := gpg_key_verifier.New(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to verify key for provider %s (cannot construct GPG key verifier: %w)", providerAddr, err)
	}

	provider, err := pkv.dataAPI.GetProvider(ctx, providerAddr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s (%w)", providerAddr, err)
	}

	toCheck := min(uint8(len(provider.Versions)), pkv.versionsToCheck)
	var matchedVersions []string

	for _, version := range provider.Versions[:toCheck] {
		shaSumContents, err := pkv.downloadFile(ctx, version.SHASumsURL)
		if err != nil {
			pkv.logger.Error("failed to download SHASums URL for provider", slog.String("provider", providerAddr.String()), slog.String("version", string(version.Version)), slog.Any("err", err))
		}

		shaSumSigContents, err := pkv.downloadFile(ctx, version.SHASumsSignatureURL)
		if err != nil {
			pkv.logger.Error("failed to download SHASums signature URL for provider", slog.String("provider", providerAddr.String()), slog.String("version", string(version.Version)), slog.Any("err", err))
		}

		if err := gpgVerifier.ValidateSignature(shaSumContents, shaSumSigContents); err != nil {
			pkv.logger.Error("failed to validate signature for provider", slog.String("provider", providerAddr.String()), slog.String("version", string(version.Version)), slog.Any("err", err))
		}

		matchedVersions = append(matchedVersions, string(version.Version))
	}

	return matchedVersions, nil
}
