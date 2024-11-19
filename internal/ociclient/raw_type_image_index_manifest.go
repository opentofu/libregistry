// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// OCIRawImageIndexManifest describes an index of OCIRawImageManifest documents. You can use
// this to determine which manifest applies to your current architecture.
// This structure conforms to the application/vnd.oci.image.index.v1+json and
// application/vnd.docker.distribution.manifest.list.v2+json media types.
//
// For details see https://github.com/opencontainers/image-spec/blob/v1.0.1/image-index.md and
// https://github.com/openshift/docker-distribution/blob/master/docs/spec/manifest-v2-2.md#manifest-list
type OCIRawImageIndexManifest struct {
	SchemaVersion int             `json:"schemaVersion"`
	MediaType     OCIRawMediaType `json:"mediaType"`
	Manifests     []struct {
		OCIRawDescriptor
		Platform struct {
			// Architecture describes the OS architecture in GOARCH strings.
			Architecture string `json:"architecture"`
			// OS describes the operating system in GOOS strings.
			OS string `json:"os"`
			// OSVersion describes the specific operating system version required.
			OSVersion  string   `json:"os.version,omitempty"`
			OSFeatures []string `json:"os.features,omitempty"`
			Variant    string   `json:"variant,omitempty"`
			Features   []string `json:"features,omitempty"`
		} `json:"platform,omitempty"`
	} `json:"manifests"`
	Annotations OCIRawAnnotations `json:"annotations,omitempty"`
}

func (o OCIRawImageIndexManifest) GetMediaType() OCIRawMediaType {
	return o.MediaType
}

var _ OCIRawManifest = OCIRawImageIndexManifest{}
