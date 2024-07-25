// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opentofu/libregistry/types/module"
)

func (r registryDataAPI) PutModule(ctx context.Context, moduleAddr module.Addr, metadata module.Metadata) error {
	marshalled, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal module metadata (%w)", err)
	}
	path := r.getModulePath(moduleAddr)
	if err := r.storageAPI.PutFile(ctx, path, marshalled); err != nil {
		return fmt.Errorf("failed to write module file %s (%w)", path, err)
	}
	return nil
}
