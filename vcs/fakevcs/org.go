// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package fakevcs

import (
	"github.com/opentofu/libregistry/vcs"
)

type org struct {
	users        map[string]struct{}
	repositories map[vcs.RepositoryAddr]*repository
}
