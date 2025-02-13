// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkeyverifier

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderDownload(t *testing.T) {
	t.Parallel()

	expectedData := []byte("test")
	key := generateKey(t)
	srv := newTestServer(t, key, expectedData)

	pkv := setupProviderCall(t, "/", "/")

	ctx := context.Background()
	data, err := pkv.(*providerKeyVerifier).downloadFile(ctx, srv.URL)

	require.NoError(t, err)

	assert.Equal(t, data, expectedData)

}
