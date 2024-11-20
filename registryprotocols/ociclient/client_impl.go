// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"context"
	"github.com/opentofu/libregistry/logger"
)

type ociClient struct {
	tempDirectory string
	rawClient     RawOCIClient
	logger        logger.Logger
}

func (o ociClient) ListReferences(ctx context.Context, addr OCIAddr) ([]OCIReference, OCIWarnings, error) {
	if err := addr.Validate(); err != nil {
		return nil, nil, err
	}
	o.logger.Debug(ctx, "Listing references for OCI image %s...", addr)
	warnings, err := o.rawClient.Check(ctx, addr.Registry)
	if err != nil {
		return nil, warnings, err
	}
	response, newWarnings, err := o.rawClient.ContentDiscovery(ctx, addr)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}
	return response.Tags, warnings, err
}

func (o ociClient) PullImage(ctx context.Context, addrRef OCIAddrWithReference, opts ...ClientPullOpt) (PulledOCIImage, OCIWarnings, error) {
	//TODO implement me
	panic("implement me")
}

var _ OCIClient = &ociClient{}
