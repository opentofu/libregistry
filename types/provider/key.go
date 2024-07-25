// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider

// Key is the key data with a key ID.
type Key struct {
	ASCIIArmor string `json:"ascii_armor"`
	KeyID      string `json:"key_id"`
}
