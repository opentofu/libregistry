// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package fakevcs

import (
	"io/fs"

	"github.com/opentofu/libregistry/vcs"
)

type version struct {
	name     vcs.Version
	assets   map[vcs.AssetName][]byte
	contents fs.ReadDirFS
}
