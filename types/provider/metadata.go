// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider

// Metadata contains information about the provider.
type Metadata struct {
	CustomRepository string    `json:"repository,omitempty"` // Optional. Custom repository from which to fetch the provider's metadata.
	Versions         []Version `json:"versions"`             // A list of version data, for each supported provider version.
}

func (m Metadata) Equals(other Metadata) bool {
	if m.CustomRepository != other.CustomRepository {
		return false
	}
	if len(m.Versions) != len(other.Versions) {
		return false
	}
	for i, version := range m.Versions {
		if !version.Equals(other.Versions[i]) {
			return false
		}
	}
	return true
}
