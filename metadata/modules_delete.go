// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"

	"github.com/opentofu/libregistry/types/module"
)

func (r registryDataAPI) DeleteModule(ctx context.Context, moduleAddr module.Addr) error {
	return r.storageAPI.DeleteFile(ctx, r.getModulePath(moduleAddr))
}
