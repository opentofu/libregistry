// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"fmt"

	"github.com/opentofu/libregistry/vcs"
	regaddr "github.com/opentofu/registry-address"
)

// Addr describes a module address combination of NAMESPACE-NAME-TARGETSYSTEM. This will translate to
// github.com/NAMESPACE/terraform-TARGETSYSTEM-NAME for now.
type Addr struct {
	Namespace    string
	Name         string
	TargetSystem string
}

func (a Addr) Validate() error {
	_, err := regaddr.ParseModuleSource(fmt.Sprintf("%s/%s/%s", a.Namespace, a.Name, a.TargetSystem))
	if err != nil {
		return &InvalidModuleAddrError{
			a,
			err,
		}
	}
	return nil
}

func (a Addr) Normalize() Addr {
	return NormalizeAddr(a)
}

func (a Addr) String() string {
	normalized := a.Normalize()
	return normalized.Namespace + "/terraform-" + normalized.Name + "-" + normalized.TargetSystem
}

func (a Addr) Equals(other Addr) bool {
	current := a.Normalize()
	other = other.Normalize()
	return current.Namespace == other.Namespace && current.Name == other.Namespace && current.TargetSystem == other.TargetSystem
}

func (a Addr) ToRepositoryAddr() vcs.RepositoryAddr {
	return vcs.RepositoryAddr{
		Org:  vcs.OrganizationAddr(a.Namespace),
		Name: "terraform-" + a.TargetSystem + "-" + a.Name,
	}
}
