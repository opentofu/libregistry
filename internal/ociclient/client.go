// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"context"
)

// OCIClient implements a high level interface giving you access to a Docker / OCI Registry.
type OCIClient interface {
	// ListVersions lists all available version references.
	ListVersions(
		ctx context.Context,
		addr OCIAddr,
	) ([]OCIReference, OCIWarnings, error)

	// PullImage implements pulling an image and returning its contents. Make sure you call Close() on the returned
	// image when you are done using it as this will clean up the temporary files.
	PullImage(
		ctx context.Context,
		addr OCIAddrWithReference,
		opts ...ClientPullOpt,
	) (PulledOCIImage, OCIWarnings, error)
}

// New creates a new OCIClient instance with the given options.
func New(opts ...Opt) (OCIClient, error) {
	c := Config{}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	return NewFromConfig(c)
}

// NewFromConfig creates a new OCIClient instance with the given configuration. For an easier syntax to use call New().
func NewFromConfig(config Config) (OCIClient, error) {
	if err := config.ApplyDefaultsAndValidate(); err != nil {
		return nil, err
	}
	return &ociClient{
		config,
	}, nil
}
