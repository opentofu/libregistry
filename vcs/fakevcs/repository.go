// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package fakevcs

import (
	"github.com/opentofu/libregistry/vcs"
)

type repository struct {
	versions []version
	info     vcs.RepositoryInfo
}
