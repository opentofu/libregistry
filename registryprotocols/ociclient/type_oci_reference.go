// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

// OCIReference is a reference to a digest or tag in an OCI registry.
// For details, see https://github.com/opencontainers/distribution-spec/blob/main/spec.md
type OCIReference string

func (o OCIReference) Validate() error {
	if !ociDigestRe.MatchString(string(o)) && !ociTagRe.MatchString(string(o)) {
		return newInvalidOCIReferenceError(o, "must match "+ociDigestRe.String()+" or "+ociTagRe.String())
	}
	return nil
}

func (o OCIReference) IsDigest() bool {
	return ociDigestRe.MatchString(string(o))
}

func (o OCIReference) IsTag() bool {
	return ociTagRe.MatchString(string(o))
}

func (o OCIReference) AsDigest() (OCIDigest, bool) {
	return OCIDigest(o), o.IsDigest()
}

func (o OCIReference) AsTag() (OCITag, bool) {
	return OCITag(o), o.IsTag()
}

func (o OCIReference) Equals(other OCIReference) bool {
	return other == o
}

var _ validatable = OCIReference("")
