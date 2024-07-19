// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package provider

// Target contains information about a specific provider version for a specific target platform.
type Target struct {
	OS          string `json:"os"`           // The operating system for which the provider is built.
	Arch        string `json:"arch"`         // The architecture for which the provider is built.
	Filename    string `json:"filename"`     // The filename of the provider binary.
	DownloadURL string `json:"download_url"` // The direct URL to download the provider binary.
	SHASum      string `json:"shasum"`       // The SHA checksum of the provider binary.
}
