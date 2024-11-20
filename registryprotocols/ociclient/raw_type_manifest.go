// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

type OCIRawManifest interface {
	GetMediaType() OCIRawMediaType
	AsIndexManifest() (OCIRawImageIndexManifest, bool)
	AsImageManifest() (OCIRawImageManifest, bool)
}
