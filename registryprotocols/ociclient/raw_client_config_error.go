// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

// RawConfigurationError indicates an error in setting up the RawOCIClient. This is typically a problem with the
// user-supplied configuration or a bug in the calling code.
type RawConfigurationError struct {
	Message string
	Cause   error
}

func (c RawConfigurationError) Error() string {
	if c.Cause != nil {
		return c.Message + " (" + c.Cause.Error() + ")"
	}
	return c.Message
}

func (c RawConfigurationError) Unwrap() error {
	return c.Cause
}

func newRawConfigurationError(message string, cause error) error { // nolint:unused // This is needed if RawConfigurationError will be used in the future.
	return &RawConfigurationError{
		Message: message,
		Cause:   cause,
	}
}
