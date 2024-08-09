// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package module

// Version represents a single version of a module.
//
// swagger:model ModuleVersion
type Version struct {
	// Version number of the provider. Correlates to a tag in the module repository.
	Version VersionNumber `json:"version"`
}

func (v Version) Normalize() Version {
	return Version{
		Version: v.Version.Normalize(),
	}
}

func (v Version) Equals(other Version) bool {
	return v.Version.Normalize() == other.Version.Normalize()
}

func (v Version) Compare(other Version) int {
	return v.Version.Compare(other.Version)
}

func (v Version) Validate() error {
	return v.Version.Validate()
}
