// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"

	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) DeleteProvider(ctx context.Context, providerAddr provider.Addr) error {
	return r.storageAPI.DeleteFile(ctx, r.getProviderPathRaw(providerAddr))
}
