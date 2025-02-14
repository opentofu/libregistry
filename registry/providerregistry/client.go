// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerregistry

import (
	"context"
)

type Client interface {
	GetProviderVersion(
		ctx context.Context,
		request GetProviderVersionRequest,
	) (GetProviderVersionResponse, error)

	ListAvailableProviderVersions(
		ctx context.Context,
		request ListAvailableProviderVersionsRequest,
	) (ListAvailableProviderVersionsResponse, error)
}
