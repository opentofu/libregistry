// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// PullConfigurationError indicates an error in setting up the OCIClient. This is typically a problem with the
// user-supplied configuration or a bug in the calling code.
type PullConfigurationError struct {
	Message string
	Cause   error
}

func (c PullConfigurationError) Error() string {
	if c.Cause != nil {
		return c.Message + " (" + c.Cause.Error() + ")"
	}
	return c.Message
}

func (c PullConfigurationError) Unwrap() error {
	return c.Cause
}

func newPullConfigurationError(message string, cause error) error {
	return &PullConfigurationError{
		Message: message,
		Cause:   cause,
	}
}
