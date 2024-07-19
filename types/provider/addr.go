// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package provider

// Addr represents a full provider address (NAMESPACE/NAME). It currently translates to
// github.com/NAMESPACE/terraform-provider-NAME .
type Addr struct {
	Namespace string
	Name      string
}

func (a Addr) Normalize() Addr {
	return NormalizeAddr(a)
}

func (a Addr) String() string {
	normalized := a.Normalize()
	return normalized.Namespace + "/terraform-provider-" + normalized.Name
}

func (a Addr) Equals(other Addr) bool {
	normalizedA := a.Normalize()
	normalizedOther := other.Normalize()
	return normalizedA.Namespace == normalizedOther.Namespace && normalizedA.Name == normalizedOther.Name
}
