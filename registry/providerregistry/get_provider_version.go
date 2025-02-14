// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerregistry

import (
	"github.com/opentofu/libregistry/registry/common"
)

// @formatter:off

// swagger:operation GET /{namespace}/{type}/{version}/download/{os}/{arch} GetProviderVersion
//
// Get a provider version
//
// Returns the details on how to download a provider for a specific OS and architecture.
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
//     description: Namespace of the provider.
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
//   - name: version
//     in: path
//     required: true
//     description: Version number of the provider
//     example: aws
//     schema:
//       '$ref': '#/definitions/RegistryProviderVersion'
//   - name: os
//     in: path
//     required: true
//     description: Operating system for the provider.
//     example: darwin
//     schema:
//       '$ref': '#/definitions/RegistryProviderOS'
//   - name: arch
//     in: path
//     required: true
//     description: Architecture for the provider.
//     example: amd64
//     schema:
//       '$ref': '#/definitions/RegistryProviderArch'
// responses:
//   '200':
//     description: The details on the version.
//     schema:
//       '$ref': '#/definitions/RegistryGetProviderVersionResponse'

// @formatter:on

type GetProviderVersionRequest struct {
	Namespace common.ProviderNamespace
	Type      common.ProviderType
	Version   common.ProviderVersion
	OS        common.ProviderOS
	Arch      common.ProviderArch
}

// GetProviderVersionResponse is the response to a GetProviderVersionRequest to obtain the download
// details of a specific provider version.
//
// swagger:model RegistryGetProviderVersionResponse
type GetProviderVersionResponse struct {
	// required: true
	Protocols []common.ProviderProtocolVersion `json:"protocols"`
	// required: true
	OS common.ProviderOS `json:"os"`
	// required: true
	Arch common.ProviderArch `json:"arch"`
	// required: true
	Filename string `json:"filename"`
	// required: true
	DownloadURL string `json:"download_url"`
	// required: true
	SHASumsURL string `json:"shasums_url"`
	// required: true
	SHASumsSignatureURL string `json:"shasums_signature_url"`
	// required: true
	SHASum string `json:"shasum"`
	// required: true
	SigningKeys struct {
		// required: true
		GPGPublicKeys []struct {
			// PGP key ID for the signing key for this provider. May be empty for the OpenTofu registry as not all
			// providers have submitted a signing key.
			// required: false
			KeyID string `json:"key_id"`
			// PGP signing key for the provider. May be empty for the OpenTofu registry as not all providers have
			// submitted a signing key.
			// required: false
			ASCIIArmor string `json:"ascii_armor"`
			// Legacy Hashicorp partner signature (PGP ASCII armor).
			// required: false
			TrustSignature string `json:"trust_signature,omitempty"`
			// Source description for the provider author. (Not used in OpenTofu or legacy Terraform.)
			Source string `json:"source,omitempty"`
			// Source URL for the provider author. (Not used in OpenTofu or legacy Terraform.)
			SourceURL string `json:"source_url,omitempty"`
		} `json:"gpg_public_keys"`
	} `json:"signing_keys"`
}
