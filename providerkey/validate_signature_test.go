// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/types/provider"
	"github.com/stretchr/testify/require"
)

func TestValidSignature(t *testing.T) {
	t.Parallel()

	key := generateKey(t)
	pubKey := getPubKey(t, key)
	signature, data := generateSignedData(t, key, []byte("test\n"))

	pk, err := New(pubKey, nil)
	require.NoError(t, err)

	p := provider.Addr{Name: "test"}
	err = pk.ValidateSignature(context.Background(), p, signature, data)
	require.NoError(t, err)
}

func TestInvalidSignature(t *testing.T) {
	t.Parallel()

	key1 := generateKey(t)
	signature, data := generateSignedData(t, key1, []byte("test\n"))

	key2 := generateKey(t)
	pubKey2 := getPubKey(t, key2)
	pk, err := New(pubKey2, nil)
	require.NoError(t, err)

	p := provider.Addr{Name: "test"}
	err = pk.ValidateSignature(context.Background(), p, signature, data)
	require.Error(t, err)
}
