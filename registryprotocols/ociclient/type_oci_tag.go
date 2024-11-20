// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"regexp"
)

// OCITag describes a tag name for OCI manifests.
// For details see https://github.com/opencontainers/image-spec/blob/v1.0.1/descriptor.md#digests
type OCITag OCIReference

func (o OCITag) Validate() error {
	if !ociTagRe.MatchString(string(o)) {
		return newInvalidOCITagError(o, "must match "+ociTagRe.String())
	}
	return nil
}

func (o OCITag) Equals(digest OCITag) bool {
	return string(o) == string(digest)
}

var _ validatable = OCITag("")

// ociTagRe is a regular expression describing OCI tags. For more details see
// https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests
var ociTagRe = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}$`)
