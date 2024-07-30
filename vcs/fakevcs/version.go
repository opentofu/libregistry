// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package fakevcs

import (
	"io/fs"
	"time"

	"github.com/opentofu/libregistry/vcs"
)

type version struct {
	name     vcs.VersionNumber
	created  time.Time
	assets   map[vcs.AssetName][]byte
	contents fs.ReadDirFS
}
