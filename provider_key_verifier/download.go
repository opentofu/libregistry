// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key_verifier

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (pkv providerKeyVerifier) downloadFile(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request (%w)", err)
	}

	response, err := pkv.config.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("status code different from 200 %s: %w", url, err)
	}

	contents := new(strings.Builder)
	_, err = io.Copy(contents, response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read file from %s: %w", url, err)
	}

	return contents.String(), nil
}
