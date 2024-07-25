// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"context"

	"github.com/opentofu/libregistry/types/module"
)

func (r registryDataAPI) GetAllModules(ctx context.Context) (map[module.Addr]module.Metadata, error) {
	moduleAddrs, err := r.ListModules(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[module.Addr]module.Metadata, len(moduleAddrs))
	for _, moduleAddr := range moduleAddrs {
		result[moduleAddr], err = r.GetModule(ctx, moduleAddr)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
