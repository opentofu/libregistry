// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkeyverifier

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/opentofu/libregistry/internal/retry"
)

func (pkv *providerKeyVerifier) downloadFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request (%w)", err)
	}

	var res *http.Response
	if err := retry.Func(
		ctx,
		fmt.Sprintf("retry file download: %s", url),
		func() error {
			res, err = pkv.config.HTTPClient.Do(req)
			return err
		},
		func(err error) bool {
			// All errors returned by HTTPClient.Do will be subject to trying to download the file again
			return true
		},
		3,
		400*time.Millisecond,
		pkv.config.Logger,
	); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code different from %d %s: %w", http.StatusOK, url, err)
	}

	contents, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from %s: %w", url, err)
	}

	return contents, nil
}
