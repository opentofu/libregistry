// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"github.com/opentofu/libregistry/types/provider"
)

type ProviderNotFoundError struct {
	ProviderAddr provider.Addr
	Cause        error
}

func (m ProviderNotFoundError) Error() string {
	return "Provider not found: " + m.ProviderAddr.String()
}

func (m ProviderNotFoundError) Unwrap() error {
	return m.Cause
}
