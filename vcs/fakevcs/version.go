// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package fakevcs

import (
	"github.com/opentofu/libregistry/vcs"
)

type version struct {
	name   vcs.Version
	assets map[vcs.AssetName][]byte
}
