// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/opentofu/libregistry/vcs"
)

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

func (a Addr) ToRepositoryAddr() vcs.RepositoryAddr {
	return vcs.RepositoryAddr{
		Org:  vcs.OrganizationAddr(a.Namespace),
		Name: "terraform-provider-" + a.Name,
	}
}
