// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

// OCIRawDescriptor describes the generic content descriptor data structure.
// For details see https://github.com/opencontainers/image-spec/blob/v1.0.1/descriptor.md
type OCIRawDescriptor struct {
	MediaType   OCIRawMediaType   `json:"mediaType"`
	Digest      OCIDigest         `json:"digest"`
	Size        int64             `json:"size"`
	URLs        []string          `json:"urls,omitempty"`
	Annotations OCIRawAnnotations `json:"annotations,omitempty"`
}
