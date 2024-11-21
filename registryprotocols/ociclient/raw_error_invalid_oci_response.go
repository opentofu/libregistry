// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

type InvalidOCIResponseError struct {
	Message string
	Cause   error
}

func (i InvalidOCIResponseError) Error() string {
	if i.Cause != nil {
		if i.Message != "" {
			return "Invalid OCI response: " + i.Message + " (" + i.Cause.Error() + ")"
		}
		return "Invalid OCI response (" + i.Cause.Error() + ")"
	}
	if i.Message != "" {
		return "Invalid OCI response: " + i.Message
	}
	return "Invalid OCI response"
}

func (i InvalidOCIResponseError) Unwrap() error {
	return i.Cause
}

func newInvalidOCIResponseError(message string, cause error) error {
	return &InvalidOCIResponseError{
		Message: message, Cause: cause,
	}
}
