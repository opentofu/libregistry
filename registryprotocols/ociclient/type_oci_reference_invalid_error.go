// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

// InvalidOCIReferenceError indicates that an invalid reference for an OCI artifact was encountered.
// This error typically happens when the user entered an invalid repository reference.
//
// You can check for this error with the following code:
//
//	var invalidOCIReferenceError *ociclient.InvalidOCIReferenceError
//	if errors.As(err, &invalidOCIReferenceError) {
//	    // Do something here
//	}
//
// Note: this error always comes wrapped in an *ociclient.ValidationError. You can use this to check
// for all validation errors.
//
// Warning: always use errors.As instead of type assertions to check for this error.
type InvalidOCIReferenceError struct {
	OCIReference OCIReference
	Reason       string
}

// Error returns a human-readable error message.
func (i InvalidOCIReferenceError) Error() string {
	if i.Reason != "" {
		return "Invalid OCI reference: " + string(i.OCIReference) + " (" + i.Reason + ")"
	}
	return "Invalid OCI reference: " + string(i.OCIReference)
}

// newInvalidOCIReferenceError provides a standardized way to construct an InvalidOCIReferenceError with the proper
// wrapping.
func newInvalidOCIReferenceError(reference OCIReference, reason string) error {
	return newValidationError(
		&InvalidOCIReferenceError{
			OCIReference: reference,
			Reason:       reason,
		},
	)
}

var _ error = &InvalidOCIReferenceError{}
