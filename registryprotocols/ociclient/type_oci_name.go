// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "regexp"

// OCIName is a name of an OCI image. Note that the name may contain a /.
// For details, see https://github.com/opencontainers/distribution-spec/blob/main/spec.md
type OCIName string

func (o OCIName) Validate() error {
	if !ociNameRe.MatchString(string(o)) {
		return newInvalidOCINameError(o, "must match "+ociNameRe.String())
	}
	return nil
}

func (o OCIName) Equals(name OCIName) bool {
	return o == name
}

var _ validatable = OCIName("")

// ociNameRe encodes the naming rules for OCI names. Note that the name may contain one or more / characters.
// For details, see https://github.com/opencontainers/distribution-spec/blob/main/spec.md
var ociNameRe = regexp.MustCompile(`^[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*(/[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*)*$`)
