// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

// OCIImageMetadata describes a minimal amount of information extracted from an OCI image.
type OCIImageMetadata struct {
	Addr   OCIAddrWithReference `json:"addr"`
	Layers []OCIDigest          `json:"layers"`
}
