// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpgvalidator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidSignature(t *testing.T) {
	t.Parallel()

	key := generateKey(t)
	signature, data := generateSignedData(t, key, []byte("test\n"))

	pk, err := New(key)
	require.NoError(t, err)

	err = pk.ValidateSignature(context.Background(), signature, data)
	require.NoError(t, err)
}

func TestInvalidSignature(t *testing.T) {
	t.Parallel()

	key1 := generateKey(t)
	signature, data := generateSignedData(t, key1, []byte("test\n"))

	key2 := generateKey(t)
	pk, err := New(key2)
	require.NoError(t, err)

	err = pk.ValidateSignature(context.Background(), signature, data)
	require.Error(t, err)
}
