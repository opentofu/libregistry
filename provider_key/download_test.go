// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"context"
	"testing"
)

func TestProviderDownload(t *testing.T) {
	expectedData := "test"
	srv := generateTestServer(t, expectedData)
	httpClient := srv.Client()

	pubKey := generateTestPubKey(t)

	pkv, err := New(pubKey, nil, WithHTTPClient(httpClient))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	ctx := context.Background()
	data, err := pkv.(*providerKey).downloadFile(ctx, srv.URL)

	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	if data != expectedData {
		t.Fatalf("expected file data is: %s, got %s instead", expectedData, data)
	}

}
