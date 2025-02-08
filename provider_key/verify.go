// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/opentofu/libregistry/types/provider"
)

type validationError struct {
	message error
}

func (e *validationError) Error() string {
	return fmt.Sprintf("%s", e.message)
}

func (p *providerKey) VerifyProvider(ctx context.Context, providerAddr provider.Addr) ([]provider.Version, error) {
	providerData, err := p.dataAPI.GetProvider(ctx, providerAddr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s (%w)", providerAddr, err)
	}

	var vError *validationError

	toCheck := min(len(providerData.Versions), int(p.config.NumVersionsToCheck))
	var matchedVersions []provider.Version

	lock := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	versions := providerData.Versions[:toCheck]
	wg.Add(len(versions))
	parallelismSemaphore := make(chan struct{}, p.config.MaxParallelism)
	errorsCh := make(chan error, len(versions))

	for _, version := range providerData.Versions[:toCheck] {
		version := version
		go func() error {

			parallelismSemaphore <- struct{}{}
			defer func() {
				<-parallelismSemaphore
			}()
			defer wg.Done()
			if err := p.check(ctx, version); err != nil {
				// If the error is different from validation, we return the error.
				// If validation is failing, func is still returning because we still want the matched versions
				errorsCh <- err
				if !errors.As(err, &vError) {
					errorsCh <- err
					return err
				}
			}

			lock.Lock()
			matchedVersions = append(matchedVersions, version)
			lock.Unlock()

			return nil
		}()
	}

	wg.Wait()
	select {
	case err := <-errorsCh:
		return nil, err
	default:
		return matchedVersions, nil
	}
}

func (p *providerKey) check(ctx context.Context, version provider.Version) error {
	shaSumContents, err := p.downloadFile(ctx, version.SHASumsURL)
	if err != nil {
		return fmt.Errorf("failed to download SHASums URL: %w", err)
	}

	signature, err := p.downloadFile(ctx, version.SHASumsSignatureURL)
	if err != nil {
		return fmt.Errorf("failed to download SHASums signature URL for provider: %w", err)
	}

	if err := p.gpgVerifier.Validate(signature, shaSumContents); err != nil {
		return &validationError{
			message: fmt.Errorf("failed to validate signature for provider: %w", err),
		}
	}
	return nil
}
