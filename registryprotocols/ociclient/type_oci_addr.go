// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "fmt"

// OCIAddr describes an address in an OCI registry.
type OCIAddr struct {
	Registry OCIRegistry `json:"registry"`
	Name     OCIName     `json:"name"`
}

func (o OCIAddr) Validate() error {
	if o.Registry == "" {
		return fmt.Errorf("invalid OCI registry address: %s (no registry specified)", o.String())
	}
	if err := o.Registry.Validate(); err != nil {
		return fmt.Errorf("invalid OCI registry address: %s (invalid registry name: %w)", o.String(), err)
	}
	if o.Name == "" {
		return fmt.Errorf("invalid OCI registry address: %s (no name specified)", o.String())
	}
	if err := o.Name.Validate(); err != nil {
		return fmt.Errorf("invalid OCI registry address: %s (invalid name: %w)", o.String(), err)
	}
	return nil
}

func (o OCIAddr) ToScope() OCIScope {
	return OCIScope{
		Registry: o.Registry,
		Name:     o.Name,
	}
}

func (o OCIAddr) Equals(other OCIAddr) bool {
	return o.Registry.Equals(other.Registry) && o.Name.Equals(other.Name)
}

func (o OCIAddr) String() string {
	return string(o.Registry) + "/" + string(o.Name)
}

var _ validatable = OCIAddr{}

// OCIAddrWithReference describes an OCIAddr with an additional reference (e.g. tag).
type OCIAddrWithReference struct {
	OCIAddr
	Reference OCIReference `json:"reference"`
}

func (o OCIAddrWithReference) Validate() error {
	if err := o.OCIAddr.Validate(); err != nil {
		return err
	}
	if err := o.Reference.Validate(); err != nil {
		return err
	}
	return nil
}

func (o OCIAddrWithReference) String() string {
	return o.OCIAddr.String() + ":" + string(o.Reference)
}

func (o OCIAddrWithReference) Equals(other OCIAddrWithReference) bool {
	return o.Registry.Equals(other.Registry) && o.Name.Equals(other.Name) && o.Reference.Equals(other.Reference)
}

var _ validatable = OCIAddrWithReference{}

// OCIAddrWithDigest describes an OCIAddr with an additional digest (e.g. SHA checksum).
type OCIAddrWithDigest struct {
	OCIAddr
	Digest OCIDigest `json:"digest"`
}

func (o OCIAddrWithDigest) Validate() error {
	if err := o.OCIAddr.Validate(); err != nil {
		return err
	}
	if err := o.Digest.Validate(); err != nil {
		return err
	}
	return nil
}

func (o OCIAddrWithDigest) String() string {
	return o.OCIAddr.String() + ":" + string(o.Digest)
}

func (o OCIAddrWithDigest) Equals(other OCIAddrWithDigest) bool {
	return o.Registry.Equals(other.Registry) && o.Name.Equals(other.Name) && o.Digest.Equals(other.Digest)
}

var _ validatable = OCIAddrWithDigest{}
