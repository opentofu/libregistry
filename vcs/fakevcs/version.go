// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package fakevcs

type version struct {
	name   string
	assets map[string][]byte
}
