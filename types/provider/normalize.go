// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"
)

func NormalizeNamespace(namespace string) string {
	return strings.ToLower(namespace)
}

func NormalizeName(name string) string {
	return strings.ToLower(name)
}

func NormalizeAddr(providerAddr Addr) Addr {
	return Addr{
		Namespace: NormalizeNamespace(providerAddr.Namespace),
		Name:      NormalizeName(providerAddr.Name),
	}
}
