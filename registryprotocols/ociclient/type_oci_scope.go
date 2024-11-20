// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"fmt"
	"strings"
)

type OCIScopeString string

func (o OCIScopeString) Parse() (OCIScope, error) {
	parts := strings.SplitN(string(o), "/", 2)
	var scope OCIScope
	if len(parts) == 1 {
		scope = OCIScope{
			Registry: OCIRegistry(parts[0]),
			Name:     "",
		}
	} else {
		scope = OCIScope{
			OCIRegistry(parts[0]),
			OCIName(parts[1]),
		}
	}
	return scope, scope.Validate()
}

// OCIScope define a scope for a registry or a specific name.
type OCIScope struct {
	Registry OCIRegistry `json:"registry"`
	Name     OCIName     `json:"name,omitempty"`
}

func (o OCIScope) Validate() error {
	if o.Registry == "" {
		return fmt.Errorf("invalid OCI registry scope: %s (no registry specified)", o.String())
	}
	if err := o.Registry.Validate(); err != nil {
		return fmt.Errorf("invalid OCI registry scope: %s (invalid registry name: %w)", o.String(), err)
	}
	if o.Name != "" {
		if err := o.Name.Validate(); err != nil {
			return fmt.Errorf("invalid OCI registry scope: %s (invalid name: %w)", o.String(), err)
		}
	}
	return nil
}

func (o OCIScope) Equals(other OCIScope) bool {
	return o.Registry.Equals(other.Registry) && o.Name == other.Name
}

func (o OCIScope) ScopeString() OCIScopeString {
	return OCIScopeString(o.String())
}

func (o OCIScope) String() string {
	return string(o.Registry) + "/" + string(o.Name)
}

func (o OCIScope) MatchesRegistry(registry OCIRegistry) bool {
	return (o.Name == "" || o.Name == "*") && o.Registry.Equals(registry)
}

func (o OCIScope) MatchesAddr(addr OCIAddr, full bool) bool {
	if full {
		return (o.Name != "" && o.Name != "*") && o.Registry.Equals(addr.Registry) && o.Name.Equals(addr.Name)
	}
	return ((o.Name != "" && o.Name != "*") && o.Registry.Equals(addr.Registry) && o.Name.Equals(addr.Name)) ||
		((o.Name == "" || o.Name == "*") && o.Registry.Equals(addr.Registry))
}
