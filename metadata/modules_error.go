// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"github.com/opentofu/libregistry/types/module"
)

type ModuleNotFoundError struct {
	ModuleAddr module.Addr
	Cause      error
}

func (m ModuleNotFoundError) Error() string {
	return "Module not found: " + m.ModuleAddr.String()
}

func (m ModuleNotFoundError) Unwrap() error {
	return m.Cause
}
