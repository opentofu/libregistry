// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "os"

// Config is the configuration structure for the OCIClient.
type Config struct {
	// TempDirectory holds the temporary directory for image layer downloads. If empty, the system temporary directory
	// will be used.
	TempDirectory string

	// RawClient holds the underlying raw client. Defaults to the built-in client.
	RawClient RawOCIClient
}

func (c *Config) ApplyDefaultsAndValidate() error {
	if c.TempDirectory == "" {
		c.TempDirectory = os.TempDir()
	}
	_, err := os.Stat(c.TempDirectory)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return newConfigurationError("invalid temporary directory: "+c.TempDirectory, err)
	}
	if err := os.MkdirAll(c.TempDirectory, 0755); err != nil {
		return newConfigurationError("cannot create temporary directory: "+c.TempDirectory, err)
	}

	if c.RawClient == nil {
		c.RawClient, err = NewRawOCIClient()
		if err != nil {
			return err
		}
	}
	return nil
}

// Opt is a builder for Config.
type Opt func(c *Config) error

// WithTempDirectory sets the temporary directory to store blobs in. If left empty, this will default to the OS
// temp directory.
func WithTempDirectory(tempDirectory string) Opt {
	return func(c *Config) error {
		c.TempDirectory = tempDirectory
		return nil
	}
}

// WithRawClient sets the underlying raw client. If not set, this will default to the built-in raw client.
func WithRawClient(client RawOCIClient) Opt {
	return func(c *Config) error {
		c.RawClient = client
		return nil
	}
}
