// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package module

import (
	"strings"
)

func NormalizeNamespace(namespace string) string {
	return strings.ToLower(namespace)
}

func NormalizeName(name string) string {
	return strings.ToLower(name)
}

func NormalizeTargetSystem(name string) string {
	return strings.ToLower(name)
}

func NormalizeAddr(moduleAddr Addr) Addr {
	return Addr{
		Namespace:    NormalizeNamespace(moduleAddr.Namespace),
		Name:         NormalizeName(moduleAddr.Name),
		TargetSystem: NormalizeTargetSystem(moduleAddr.TargetSystem),
	}
}
