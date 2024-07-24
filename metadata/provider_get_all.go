// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"context"

	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) GetAllProviders(ctx context.Context, includeAliases bool) (map[provider.Addr]provider.Metadata, error) {
	providerAddrs, err := r.ListProviders(ctx, includeAliases)
	if err != nil {
		return nil, err
	}
	result := make(map[provider.Addr]provider.Metadata, len(providerAddrs))
	for _, providerAddr := range providerAddrs {
		result[providerAddr], err = r.GetProvider(ctx, providerAddr, true)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
