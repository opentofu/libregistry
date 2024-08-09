// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentofu/libregistry/vcs"
)

// Addr represents a full provider address (NAMESPACE/NAME). It currently translates to
// github.com/NAMESPACE/terraform-provider-NAME .
//
// swagger:type string
type Addr struct {
	Namespace string
	Name      string
}

func (a Addr) MarshalJSON() ([]byte, error) {
	// Note: this intentionally doesn't have a pointer receiver! Don't add one!
	normalized := a.Normalize()
	//goland:noinspection GoRedundantConversion
	return json.Marshal(string(normalized.Namespace + "/" + normalized.Name))
}

func (a *Addr) UnmarshalJSON(b []byte) error {
	var data string
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	parts := strings.Split(data, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid provider address: %s", data)
	}
	a.Namespace = parts[0]
	a.Name = parts[1]
	return nil
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

func (a Addr) Compare(other Addr) int {
	namespaceComparison := strings.Compare(a.Namespace, other.Namespace)
	if namespaceComparison != 0 {
		return namespaceComparison
	}
	return strings.Compare(a.Name, other.Name)
}

func (a Addr) ToRepositoryAddr() vcs.RepositoryAddr {
	return vcs.RepositoryAddr{
		Org:  vcs.OrganizationAddr(a.Namespace),
		Name: "terraform-provider-" + a.Name,
	}
}
