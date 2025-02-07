// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"context"
	"testing"
)

func TestProviderDownload(t *testing.T) {
	expectedData := "test"
	httpClient := generateTestClient(expectedData)
	keyData, err := generateKey()
	if err != nil {
		t.Fatalf("couldn't create key: %v", err)
	}

	ctx := context.Background()
	pkv, err := New(keyData, nil, WithHTTPClient(httpClient))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	data, err := pkv.(*providerKey).downloadFile(ctx, "https://example.com/test.txt")

	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	if data != expectedData {
		t.Fatalf("expected file data is: %s, got %s instead", expectedData, data)
	}

}
