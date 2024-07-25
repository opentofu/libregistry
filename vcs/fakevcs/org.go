// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package fakevcs

import (
	"github.com/opentofu/libregistry/vcs"
)

type org struct {
	users        map[vcs.Username]struct{}
	repositories map[vcs.RepositoryAddr]*repository
}
