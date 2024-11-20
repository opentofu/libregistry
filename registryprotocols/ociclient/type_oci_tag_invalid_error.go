// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// InvalidOCITagError indicates that an invalid tag for an OCI artifact was encountered.
// This error typically happens when the user entered an invalid repository tag.
//
// You can check for this error with the following code:
//
//	var invalidOCITagError *ociclient.InvalidOCITagError
//	if errors.As(err, &invalidOCITagError) {
//	    // Do something here
//	}
//
// Note: this error always comes wrapped in an *ociclient.ValidationError. You can use this to check
// for all validation errors.
//
// Warning: always use errors.As instead of type assertions to check for this error.
type InvalidOCITagError struct {
	OCITag OCITag
	Reason string
}

// Error returns a human-readable error message.
func (i InvalidOCITagError) Error() string {
	if i.Reason != "" {
		return "Invalid OCI tag: " + string(i.OCITag) + " (" + i.Reason + ")"
	}
	return "Invalid OCI tag: " + string(i.OCITag)
}

// newInvalidOCITagError provides a standardized way to construct an InvalidOCITagError with the proper
// wrapping.
func newInvalidOCITagError(tag OCITag, reason string) error {
	return newValidationError(
		&InvalidOCITagError{
			OCITag: tag,
			Reason: reason,
		},
	)
}

var _ error = &InvalidOCITagError{}
