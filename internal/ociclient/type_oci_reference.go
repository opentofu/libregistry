// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "regexp"

// OCIReference is a reference to a digest or tag in an OCI registry.
// For details, see https://github.com/opencontainers/distribution-spec/blob/main/spec.md
type OCIReference string

func (o OCIReference) Validate() error {
	return validateOCIReference(o)
}

func validateOCIReference(o OCIReference) error {
	if !ociReferenceRe.MatchString(string(o)) {
		return newInvalidOCIReferenceError(o, "must match "+ociReferenceRe.String())
	}
	return nil
}

func (o OCIReference) Equals(other OCIReference) bool {
	return other == o
}

var _ validatable = OCIReference("")

// ociReferenceRe encodes the reference naming rules for OCI references (e.g. digests or tags.)
// For details, see https://github.com/opencontainers/distribution-spec/blob/main/spec.md
var ociReferenceRe = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}$`)
