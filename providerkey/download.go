// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (pk *providerKey) downloadFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request (%w)", err)
	}

	res, err := pk.config.HTTPClient.Do(req)
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
