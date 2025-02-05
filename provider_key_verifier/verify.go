// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key_verifier

import (
	"context"
	"fmt"

	"github.com/opentofu/libregistry/types/provider"
)

func (pkv providerKeyVerifier) VerifyProvider(ctx context.Context, providerAddr provider.Addr) ([]*provider.Version, error) {
	providerData, err := pkv.dataAPI.GetProvider(ctx, providerAddr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s (%w)", providerAddr, err)
	}

	toCheck := min(len(providerData.Versions), int(pkv.config.NumVersionsToCheck))
	matchedVersions := make([]*provider.Version, 0)
	versionChan := make(chan *provider.Version)

	for _, version := range providerData.Versions[:toCheck] {
		go func(version provider.Version) {
			err := pkv.config.checkFn(pkv, ctx, version)
			if err != nil {
				// pkv.config.Logger.Error("error in version:", slog.String("provider", providerAddr.String()), slog.String("version", string(version.Version)), slog.Any("err", err))
				versionChan <- nil
				return
			}
			versionChan <- &version

		}(version)
	}

	for i := 0; i < toCheck; i++ {
		v := <-versionChan
		if v != nil {
			matchedVersions = append(matchedVersions, v)
		}
	}

	return matchedVersions, nil
}

func process(pkv providerKeyVerifier, ctx context.Context, version provider.Version) error {
	shaSumContents, err := pkv.downloadFile(ctx, version.SHASumsURL)
	if err != nil {
		return fmt.Errorf("failed to download SHASums URL")
	}

	shaSumSigContents, err := pkv.downloadFile(ctx, version.SHASumsSignatureURL)
	if err != nil {
		return fmt.Errorf("failed to download SHASums signature URL for provider")
	}

	if err := pkv.gpgVerifier.ValidateSignature(shaSumContents, shaSumSigContents); err != nil {
		return fmt.Errorf("failed to validate signature for provider")
	}
	return nil
}
