// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderConfig(t *testing.T) {
	t.Parallel()

	key := generateKey(t)
	pubKey := getPubKey(t, key)

	versionsToCheck := uint8(5)
	pkv, err := New(pubKey, nil, WithNumVersionsToCheck(versionsToCheck))

	require.NoError(t, err)

	assert.Equal(t, pkv.(*providerKey).config.VersionsToCheck, versionsToCheck)
}
