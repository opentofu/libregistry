// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providermirrorprotocol

import (
	"github.com/opentofu/libregistry/registry/common"
)

type ListAvailableInstallationPackagesRequest struct {
	Hostname  common.Hostname
	Namespace common.ProviderNamespace
	Name      common.ProviderType
	Version   common.ProviderVersion
}

// swagger:model ListAvailableInstallationPackagesResponse
type ListAvailableInstallationPackagesResponse struct {
	// A map of platform/architecture combinations to the details
	Archives map[common.PackageTarget]ListAvailableInstallationPackagesArchive `json:"archives"`
}

// ListAvailableInstallationPackagesArchive is a single item in a ListAvailableInstallationPackagesResponse.
//
// swagger:model ListAvailableInstallationPackagesArchive
type ListAvailableInstallationPackagesArchive struct {
	// URL to a ZIP archive containing the installation package. This can be relative
	// to the current URL.
	//
	// required: true
	URL string `json:"url"`
	// A list of hashes for the ZIP archive.
	//
	// TODO: describe the hash format.
	//
	// required: true
	Hashes []ListAvailableInstallationPackagesHash `json:"hashes"`
}

// ListAvailableInstallationPackagesHash is a single hash in an ListAvailableInstallationPackagesArchive.
//
// min length: 1
// swagger:model ListAvailableInstallationPackagesHash
type ListAvailableInstallationPackagesHash string
