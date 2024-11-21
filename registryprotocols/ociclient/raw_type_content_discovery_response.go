// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

type OCIRawContentDiscoveryResponse struct {
	Name OCIName        `json:"name"`
	Tags []OCIReference `json:"tags"`
}
