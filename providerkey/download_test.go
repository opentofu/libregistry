// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderDownload(t *testing.T) {
	expectedData := []byte("test")
	key := generateKey(t)
	srv := newTestServer(t, key, expectedData)
	httpClient := srv.Client()

	pubKey := getPubKey(t, key)
	pkv, err := New(pubKey, nil, WithHTTPClient(httpClient))

	require.NoError(t, err)

	ctx := context.Background()
	data, err := pkv.(*providerKey).downloadFile(ctx, srv.URL)

	require.NoError(t, err)

	assert.Equal(t, data, expectedData)

}
