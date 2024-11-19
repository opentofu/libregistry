// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"fmt"
	"strings"
)

// OCIRawErrors is an error response returned from an OCI registry API.
type OCIRawErrors struct {
	// TODO: this looks similar to hcl.Diagnostics. Do we want to support converting it into one?
	// TODO: this may not be ideal as then libregistry would depend on HCL.

	Errors []OCIRawError `json:"errors"`
}

// Has returns true if at least one error has the specified error code.
func (o OCIRawErrors) Has(code OCIRawErrorCode) bool {
	for _, e := range o.Errors {
		if e.Code == code {
			return true
		}
	}
	return false
}

// Errors creates a compound error string.
func (o OCIRawErrors) Error() string {
	errors := make([]string, len(o.Errors))
	for i, err := range o.Errors {
		errors[i] = err.Error()
	}
	return fmt.Sprintf("One or more OCI errors have occured (%s)", strings.Join(errors, "; "))
}

// OCIRawError is a single error in an OCIRawErrors response.
type OCIRawError struct {
	Code    OCIRawErrorCode `json:"code"`
	Message string          `json:"message"`
	Detail  string          `json:"detail"`
}

// OCIRawError attempts to create a human-readable string from the error.
func (o OCIRawError) Error() string {
	if o.Detail != "" {
		return fmt.Sprintf("%s %s (%s)", o.Code, o.Message, o.Detail)
	}
	return fmt.Sprintf("%s %s", o.Code, o.Message)
}
