// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// ConfigurationError indicates an error in setting up the OCIClient. This is typically a problem with the
// user-supplied configuration or a bug in the calling code.
type ConfigurationError struct {
	Message string
	Cause   error
}

func (c ConfigurationError) Error() string {
	if c.Cause != nil {
		return c.Message + " (" + c.Cause.Error() + ")"
	}
	return c.Message
}

func (c ConfigurationError) Unwrap() error {
	return c.Cause
}

func newConfigurationError(message string, cause error) error {
	return &ConfigurationError{
		Message: message,
		Cause:   cause,
	}
}
