// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentofu/libregistry/vcs"
	regaddr "github.com/opentofu/registry-address"
)

// Addr describes a module address combination of NAMESPACE-NAME-TARGETSYSTEM. This will translate to
// github.com/NAMESPACE/terraform-TARGETSYSTEM-NAME for now.
//
// swagger:type string
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

func (a Addr) MarshalJSON() ([]byte, error) {
	// Note: this intentionally doesn't have a pointer receiver! Don't add one!
	normalized := a.Normalize()
	//goland:noinspection GoRedundantConversion
	return json.Marshal(string(normalized.Namespace + "/" + normalized.Name + "/" + normalized.TargetSystem))
}

func (a *Addr) UnmarshalJSON(b []byte) error {
	var data string
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	parts := strings.Split(data, "/")
	if len(parts) != 3 {
		return fmt.Errorf("invalid module address: %s", data)
	}
	a.Namespace = parts[0]
	a.Name = parts[1]
	a.TargetSystem = parts[2]
	return nil
}

func (a Addr) Compare(other Addr) int {
	namespaceComparison := strings.Compare(a.Namespace, other.Namespace)
	if namespaceComparison != 0 {
		return namespaceComparison
	}
	nameComparison := strings.Compare(a.Name, other.Name)
	if nameComparison != 0 {
		return nameComparison
	}
	return strings.Compare(a.TargetSystem, other.TargetSystem)
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
	return current.Namespace == other.Namespace && current.Name == other.Name && current.TargetSystem == other.TargetSystem
}

func (a Addr) ToRepositoryAddr() vcs.RepositoryAddr {
	return vcs.RepositoryAddr{
		Org:  vcs.OrganizationAddr(a.Namespace),
		Name: "terraform-" + a.TargetSystem + "-" + a.Name,
	}
}
