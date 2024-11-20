// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// OCIRawImageManifest describes a single OCI image with multiple layers.
// This structure conforms to the application/vnd.oci.image.manifest.v1+json
// and application/vnd.docker.distribution.manifest.v2+json media types.
//
// For details see: https://github.com/opencontainers/image-spec/blob/v1.0.1/manifest.md and
// https://github.com/openshift/docker-distribution/blob/master/docs/spec/manifest-v2-2.md#image-manifest
type OCIRawImageManifest struct {
	SchemaVersion int                `json:"schemaVersion"`
	MediaType     OCIRawMediaType    `json:"mediaType,omitempty"`
	Config        OCIRawDescriptor   `json:"config"`
	Layers        []OCIRawDescriptor `json:"layers"`
	Annotations   OCIRawAnnotations  `json:"annotations,omitempty"`
}

func (o OCIRawImageManifest) GetMediaType() OCIRawMediaType {
	return o.MediaType
}

var _ OCIRawManifest = OCIRawImageManifest{}
