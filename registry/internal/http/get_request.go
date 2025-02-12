// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetRequest[T any](ctx context.Context, httpClient *http.Client, endpoint string, urlSuffix string) (response T, err error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		strings.TrimSuffix(endpoint, "/")+"/"+urlSuffix,
		nil,
	)
	if err != nil {
		return response, fmt.Errorf("failed to construct HTTP request (%w)", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("HTTP request failed (%w)", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		return response, fmt.Errorf("invalid status code returned: %d", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&response); err != nil {
		return response, fmt.Errorf("protocol error, failed to decode response (%w)", err)
	}
	return response, nil
}
