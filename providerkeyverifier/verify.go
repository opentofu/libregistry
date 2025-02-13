// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkeyverifier

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/opentofu/libregistry/types/provider"
	"golang.org/x/sync/errgroup"
)

type validationError struct {
	message error
}

func (e *validationError) Error() string {
	return fmt.Sprintf("%s", e.message)
}

func (pk *providerKeyVerifier) VerifyProvider(ctx context.Context, pAddr provider.Addr) ([]provider.Version, error) {
	pk.config.Logger.Info(ctx, "Verifying provider %s...", pAddr.Name)
	providerData, err := pk.dataAPI.GetProvider(ctx, pAddr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s (%w)", pAddr, err)
	}

	toCheck := min(len(providerData.Versions), int(pk.config.VersionsToCheck))
	var signedVersions []provider.Version

	lock := &sync.Mutex{}
	parallelismSemaphore := make(chan struct{}, pk.config.MaxParallelism)
	g, ctx := errgroup.WithContext(ctx)

	for _, version := range providerData.Versions[:toCheck] {
		version := version
		g.Go(func() error {

			parallelismSemaphore <- struct{}{}
			defer func() {
				<-parallelismSemaphore
			}()
			if err := pk.validate(ctx, pAddr, version); err != nil {
				var vError *validationError
				// If it isn't a validation error, like a network error, we return it and fail the function.
				if !errors.As(err, &vError) {
					return err
				}
				// If signature's validation is failing, the version is not added to the signed versions list.
				return nil
			}

			lock.Lock()
			signedVersions = append(signedVersions, version)
			lock.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("error when verifying provider versions: %w", err)
	}

	return signedVersions, nil
}

// check is used to download the version's signature and data and validates it with the keyring.
func (pk *providerKeyVerifier) validate(ctx context.Context, pAddr provider.Addr, version provider.Version) error {
	pk.config.Logger.Info(ctx, "Validating signature with key %s for provider %s...", pk.key.GetHexKeyID(), pAddr.Name)
	shaSumContents, err := pk.downloadFile(ctx, version.SHASumsURL)
	if err != nil {
		return fmt.Errorf("failed to download SHASums URL: %w", err)
	}

	signature, err := pk.downloadFile(ctx, version.SHASumsSignatureURL)
	if err != nil {
		return fmt.Errorf("failed to download SHASums signature URL for provider: %w", err)
	}

	if err := pk.gpgValidator.ValidateSignature(ctx, signature, shaSumContents); err != nil {
		return &validationError{
			message: fmt.Errorf("failed to validate signature for provider: %w", err),
		}
	}
	return nil
}
