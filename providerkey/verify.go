// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

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

func (pk *providerKey) VerifyProvider(ctx context.Context, addr provider.Addr) ([]provider.Version, error) {
	providerData, err := pk.dataAPI.GetProvider(ctx, addr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s (%w)", addr, err)
	}

	toCheck := min(len(providerData.Versions), int(pk.config.NumVersionsToCheck))
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
			if err := pk.check(ctx, version); err != nil {
				// If the error is different from validation, we return the error.
				var vError *validationError
				if !errors.As(err, &vError) {
					return err
				}
				// If validation is failing, func is still returning because we still want the matched versions
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

func (pk *providerKey) check(ctx context.Context, version provider.Version) error {
	shaSumContents, err := pk.downloadFile(ctx, version.SHASumsURL)
	if err != nil {
		return fmt.Errorf("failed to download SHASums URL: %w", err)
	}

	signature, err := pk.downloadFile(ctx, version.SHASumsSignatureURL)
	if err != nil {
		return fmt.Errorf("failed to download SHASums signature URL for provider: %w", err)
	}

	if err := pk.ValidateSignature(signature, shaSumContents); err != nil {
		return &validationError{
			message: fmt.Errorf("failed to validate signature for provider: %w", err),
		}
	}
	return nil
}
