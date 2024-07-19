// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) PutProvider(ctx context.Context, providerAddr provider.Addr, metadata provider.Metadata) error {
	marshalled, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal module metadata (%w)", err)
	}
	path := r.getProviderPathRaw(providerAddr)
	if err := r.storageAPI.PutFile(ctx, path, marshalled); err != nil {
		return fmt.Errorf("failed to write module file %s (%w)", path, err)
	}
	return nil
}
