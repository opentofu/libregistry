// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"regexp"
	"strings"
)

// OCIDigest describes a checksum used in OCI manifests. The default registered
// algorithms include "sha256" and "sha512", whereas "sha256" should be used by default.
// For details see https://github.com/opencontainers/image-spec/blob/v1.0.1/descriptor.md#digests
type OCIDigest string

func (o OCIDigest) Validate() error {
	if !ociDigestRe.MatchString(string(o)) {
		return newInvalidOCIDigestError(o, "must match "+ociDigestRe.String(), nil)
	}
	return nil
}

func (o OCIDigest) Equals(digest OCIDigest) bool {
	return strings.ToLower(string(o)) == strings.ToLower(string(digest))
}

var _ validatable = OCIDigest("")

// ociDigestRe is a regular expression matching the algorithm grammar:
//
//	digest                ::= algorithm ":" encoded
//	algorithm             ::= algorithm-component (algorithm-separator algorithm-component)*
//	algorithm-component   ::= [a-z0-9]+
//	algorithm-separator   ::= [+._-]
//	encoded               ::= [a-zA-Z0-9=_-]+
var ociDigestRe = regexp.MustCompile(`^([a-z0-9]+([+._-][a-z0-9]+)*):([a-zA-Z0-9=_-]+)$`)
