// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"bytes"
	"context"
	"testing"
)

func TestProviderDownload(t *testing.T) {
	expectedData := []byte("test")
	key := generateKey(t)
	srv := generateTestServer(t, key, expectedData)
	httpClient := srv.Client()

	pubKey := getPubKey(t, key)
	pkv, err := New(pubKey, nil, WithHTTPClient(httpClient))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	ctx := context.Background()
	data, err := pkv.(*providerKey).downloadFile(ctx, srv.URL)

	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	if !bytes.Equal(data, expectedData) {
		t.Fatalf("Expected file data is: %s, got %s instead", expectedData, data)
	}

}
