// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

// InvalidOCIDigestError indicates that an invalid digest for an OCI artifact was encountered.
// This error typically happens when the user entered an invalid repository digest.
//
// You can check for this error with the following code:
//
//	var invalidOCIDigestError *ociclient.InvalidOCIDigestError
//	if errors.As(err, &invalidOCIDigestError) {
//	    // Do something here
//	}
//
// Note: this error always comes wrapped in an *ociclient.ValidationError. You can use this to check
// for all validation errors.
//
// Warning: always use errors.As instead of type assertions to check for this error.
type InvalidOCIDigestError struct {
	OCIDigest OCIDigest
	Reason    string
	Cause     error
}

// Error returns a human-readable error message.
func (i InvalidOCIDigestError) Error() string {
	if i.Cause != nil && i.Reason != "" {
		return "Invalid OCI digest: " + string(i.OCIDigest) + " (" + i.Reason + "; " + i.Cause.Error() + ")"
	}
	if i.Reason != "" {
		return "Invalid OCI digest: " + string(i.OCIDigest) + " (" + i.Reason + ")"
	}
	if i.Cause != nil {
		return "Invalid OCI digest: " + string(i.OCIDigest) + " (" + i.Cause.Error() + ")"
	}
	return "Invalid OCI digest " + string(i.OCIDigest)
}

// Unwrap returns the underlying error, if any.
func (i InvalidOCIDigestError) Unwrap() error {
	return i.Cause
}

// newInvalidOCIDigestError provides a standardized way to construct an InvalidOCIDigestError with the proper
// wrapping.
func newInvalidOCIDigestError(digest OCIDigest, reason string, cause error) error {
	return newValidationError(
		&InvalidOCIDigestError{
			OCIDigest: digest,
			Reason:    reason,
			Cause:     cause,
		},
	)
}

var _ error = &InvalidOCIDigestError{}
