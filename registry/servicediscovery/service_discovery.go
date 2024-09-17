// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package servicediscovery

type Request struct {
}

// Response contains the response structure for an OpenTofu service discovery request.
//
// swagger:model RegistryServiceDiscoveryResponse
type Response struct {
	// required: false
	// format: uri-reference
	ProvidersV1 string `json:"providers.v1,omitempty"`
	// required: false
	// format: uri-reference
	ModulesV1 string `json:"modules.v1,omitempty"`
	// Login protocol struct. This is not supported by the OpenTofu Registry, but alternative
	// registries may make use of them. See https://opentofu.org/docs/internals/login-protocol/
	// for details.
	// required: false
	LoginV1 *struct {
		// Client ID OpenTofu should use to authenticate.
		//
		// ---
		// required: true
		Client string `json:"client"`
		// Supported grant types. Currently only the authz_code grant type is supported in OpenTofu.
		//
		// ---
		// required: true
		GrantTypes []string `json:"grant_types"`
		// OAuth URL for the authorization endpoint.
		//
		// ---
		// required: true
		// format: uri
		AuthzURL string `json:"authz"`
		// OAuth URL for the token endpoint.
		//
		// ---
		// required: true
		// format: uri
		TokenURL string `json:"token"`
		// Port range from-to for temporary HTTP server to perform OAuth authentication.
		//
		// ---
		// required: false
		// minItems: 2
		// maxItems: 2
		// default: [1024,65535]
		// example: [10000, 10010]
		Ports []PortNumber `json:"ports,omitempty"`
	} `json:"login.v1,omitempty"`
}

// swagger:model RegistryPortNumber
type PortNumber int
