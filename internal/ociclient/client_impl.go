// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "context"

type ociClient struct {
	config Config
}

func (o ociClient) ListVersions(ctx context.Context, addr OCIAddr) ([]OCIReference, OCIWarnings, error) {
	warnings, err := o.config.RawClient.Check(ctx, addr.Registry)
	if err != nil {
		return nil, warnings, err
	}
	response, newWarnings, err := o.config.RawClient.ContentDiscovery(ctx, addr)
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
