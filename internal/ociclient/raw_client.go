// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"context"
)

// RawOCIClient implements a Docker / OCI Registry client interface directly giving you the ability to
// call the API calls. For most use cases use OCIClient instead
type RawOCIClient interface {
	// Check checks if the hostname implements the OCI registry protocol. If needed, you can use this to obtain basic
	// credentials.
	Check(
		ctx context.Context,
		registry OCIRegistry,
	) (OCIWarnings, error)

	// ContentDiscovery performs an OCI content discovery against the /v2/<name>/tags/list endpoint.
	//
	// See https://github.com/opencontainers/distribution-spec/blob/main/spec.md#content-discovery and
	// https://distribution.github.io/distribution/spec/api/#listing-image-tags for details.
	ContentDiscovery(
		ctx context.Context,
		addr OCIAddr,
	) (OCIRawContentDiscoveryResponse, OCIWarnings, error)

	// GetManifest returns a manifest from an OCI registry against the /v2/<name>/manifests/<reference> endpoint.
	// The response will either be an OCIRawImageIndexManifest or an OCIRawImageManifest. If the manifest type cannot
	// be determined, an *OCIRawUnknownManifestTypeError is returned.
	//
	// For details see https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests and
	// https://distribution.github.io/distribution/spec/api/#pulling-an-image-manifest
	GetManifest(
		ctx context.Context,
		addrRef OCIAddrWithReference,
	) (OCIRawManifest, OCIWarnings, error)

	// GetBlob fetches a blob from an OCI registry against the /v2/<name>/blobs/<digest> endpoint.
	//
	// For details see https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-blobs and
	// https://distribution.github.io/distribution/spec/api/#pulling-a-layer
	GetBlob(
		ctx context.Context,
		addrDigest OCIAddrWithDigest,
	) (OCIRawBlob, OCIWarnings, error)
}

func NewRawOCIClient(opts ...RawOCIClientOpt) (RawOCIClient, error) {
	c := RawOCIClientConfig{}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	if err := c.ApplyDefaultsAndValidate(); err != nil {
		return nil, err
	}
	return &rawClient{
		httpClient:  c.HTTPClient,
		credentials: c.Credentials,
	}, nil
}
