// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package provider

// Version contains information about a specific provider version.
type Version struct {
	Version             string   `json:"version"`               // The version number of the provider.
	Protocols           []string `json:"protocols"`             // The protocol versions the provider supports.
	SHASumsURL          string   `json:"shasums_url"`           // The URL to the SHA checksums file.
	SHASumsSignatureURL string   `json:"shasums_signature_url"` // The URL to the GPG signature of the SHA checksums file.
	Targets             []Target `json:"targets"`               // A list of target platforms for which this provider version is available.
}

func (v Version) Equals(other Version) bool {
	if v.Version != other.Version {
		return false
	}
	if len(v.Protocols) != len(other.Protocols) {
		return false
	}
	for i, proto := range v.Protocols {
		if proto != other.Protocols[i] {
			return false
		}
	}
	if v.SHASumsURL != other.SHASumsURL {
		return false
	}
	if v.SHASumsSignatureURL != other.SHASumsSignatureURL {
		return false
	}
	if len(v.Targets) != len(other.Targets) {
		return false
	}
	for i, target := range v.Targets {
		if !target.Equals(v.Targets[i]) {
			return false
		}
	}
	return true
}
