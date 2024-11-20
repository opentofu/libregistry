// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// InvalidOCINameError indicates that an invalid name for an OCI artifact was encountered.
// This error typically happens when the user entered an invalid repository name.
//
// You can check for this error with the following code:
//
//	var invalidOCINameError *ociclient.InvalidOCINameError
//	if errors.As(err, &invalidOCINameError) {
//	    // Do something here
//	}
//
// Note: this error always comes wrapped in an *ociclient.ValidationError. You can use this to check
// for all validation errors.
//
// Warning: always use errors.As instead of type assertions to check for this error.
type InvalidOCINameError struct {
	OCIName OCIName
	Reason  string
}

// Error returns a human-readable error message.
func (i InvalidOCINameError) Error() string {
	if i.Reason != "" {
		return "Invalid OCI name: " + string(i.OCIName) + " (" + i.Reason + ")"
	}
	return "Invalid OCI name: " + string(i.OCIName)
}

// newInvalidOCINameError provides a standardized way to construct an InvalidOCINameError with the proper
// wrapping.
func newInvalidOCINameError(name OCIName, reason string) error {
	return newValidationError(
		&InvalidOCINameError{
			OCIName: name,
			Reason:  reason,
		},
	)
}

var _ error = &InvalidOCINameError{}
