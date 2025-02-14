// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providermirror

import (
	"github.com/opentofu/libregistry/registry/common"
)

type ListAvailableVersionsRequest struct {
	Hostname  common.Hostname
	Namespace common.ProviderNamespace
	Type      common.ProviderType
}

// swagger:model ListAvailableVersionsResponse
type ListAvailableVersionsResponse struct {
	// required: true
	Versions map[common.ProviderVersion]ListAvailableVersionsVersion `json:"versions"`
}

// swagger:model ListAvailableVersionsVersion
type ListAvailableVersionsVersion struct {
}
