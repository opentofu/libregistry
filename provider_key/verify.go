// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"context"
	"fmt"
	"sync"

	"github.com/opentofu/libregistry/types/provider"
)

func (p *providerKey) VerifyProvider(ctx context.Context, providerAddr provider.Addr) ([]provider.Version, error) {
	providerData, err := p.dataAPI.GetProvider(ctx, providerAddr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s (%w)", providerAddr, err)
	}

	toCheck := min(len(providerData.Versions), int(p.config.NumVersionsToCheck))
	var matchedVersions []provider.Version

	lock := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	versions := providerData.Versions[:toCheck]
	wg.Add(len(versions))
	parallelismSemaphore := make(chan struct{}, p.config.MaxParallelism)
	for _, version := range providerData.Versions[:toCheck] {
		version := version
		go func() {
			parallelismSemaphore <- struct{}{}
			if err := p.check(ctx, version); err != nil {
				// p.config.Logger.Error(ctx, "failed to verify key for provider %s version %s (%v)", providerAddr.String(), string(version.Version), err)
				// return
			}
			lock.Lock()
			matchedVersions = append(matchedVersions, version)
			lock.Unlock()
			<-parallelismSemaphore
		}()
	}
	wg.Wait()

	return matchedVersions, nil
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
		return fmt.Errorf("failed to validate signature for provider: %w", err)
	}
	return nil
}
