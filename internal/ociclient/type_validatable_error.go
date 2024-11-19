// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// ValidationError indicates than an OCI-related entity (name, reference, response, etc) failed validation.
// You can check if this error was returned with the following code:
//
//	var validationError *ociclient.ValidationError
//	if errors.As(err, &validationError) {
//	    // Do something here
//	}
//
// Note that this type often wraps a more specific error. Check the detailed errors for more information.
type ValidationError struct {
	Cause error
}

// Error returns an error message, possibly containing the root cause of the error.
func (v ValidationError) Error() string {
	if v.Cause != nil {
		return "OCI validation error (" + v.Cause.Error() + ")"
	}
	return "OCI validation error"
}

// Unwrap returns the specialized error.
func (v ValidationError) Unwrap() error {
	return v.Cause
}

func newValidationError(cause error) error {
	return &ValidationError{cause}
}

var _ error = ValidationError{}
