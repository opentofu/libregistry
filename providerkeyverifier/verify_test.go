// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkeyverifier

import (
	"context"
	"testing"
	"time"

	"github.com/opentofu/libregistry/types/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderValidVerify(t *testing.T) {
	t.Parallel()

	pkv := setupProviderCall(t, "/SHASumsURL/", "/SHASumsSignatureURL/")
	addr := provider.Addr{
		Name:      "test",
		Namespace: "opentofu",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, err := pkv.VerifyProvider(ctx, addr)
	require.NoError(t, err)

	assert.NotEqual(t, len(data), 0)

	assert.Equal(t, string(data[0].Version), "0.2.0")
}

func TestProviderInvalidVerify(t *testing.T) {
	t.Parallel()

	pkv := setupProviderCall(t, "/invalid/", "/invalid/")
	addr := provider.Addr{
		Name:      "test",
		Namespace: "opentofu",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, err := pkv.VerifyProvider(ctx, addr)
	require.NoError(t, err)

	assert.Equal(t, len(data), 0)
}
