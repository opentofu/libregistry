// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerregistry

import (
	"github.com/opentofu/libregistry/registry/common"
)

// @formatter:off

// swagger:operation GET /{namespace}/{type}/versions ListAvailableProviderVersions
//
// List available provider versions
//
// Returns the endpoints for modules and providers.
//
// ---
// tags:
// - Registry Protocol
// produces:
// - application/json
// parameters:
//   - name: namespace
//     in: path
//     required: true
//     description: Namespace of the provider to mirror.
//     example: hashicorp
//     schema:
//       '$ref': '#/definitions/RegistryProviderNamespace'
//   - name: type
//     in: path
//     required: true
//     description: Type part of the provider address.
//     example: aws
//     schema:
//       '$ref': '#/definitions/RegistryProviderType'
// responses:
//   '200':
//     description: A list of versions for a provider.
//     schema:
//       '$ref': '#/definitions/RegistryListAvailableProviderVersionsResponse'

// @formatter:on

type ListAvailableProviderVersionsRequest struct {
	Namespace common.ProviderNamespace
	Type      common.ProviderType
}

// swagger:model RegistryListAvailableProviderVersionsResponse
type ListAvailableProviderVersionsResponse struct {
	Versions []ListAvailableProviderVersionsResponseVersion `json:"versions"`
	Warnings []string                                       `json:"warnings,omitempty"`
}

// swagger:model RegistryListAvailableProviderVersionsResponseVersion
type ListAvailableProviderVersionsResponseVersion struct {
	Version   common.ProviderVersion           `json:"version"`
	Protocols []common.ProviderProtocolVersion `json:"protocols"`
	Platforms []ListAvailableProviderVersionsPlatform
}

// swagger:model RegistryListAvailableProviderVersionsPlatform
type ListAvailableProviderVersionsPlatform struct {
	OS   common.ProviderOS   `json:"os"`
	Arch common.ProviderArch `json:"arch"`
}
