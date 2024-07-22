// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/module"
)

func (r registryDataAPI) GetModule(ctx context.Context, moduleAddr module.Addr) (module.Metadata, error) {
	path := r.getModulePath(moduleAddr)
	fileContents, err := r.storageAPI.GetFile(ctx, path)
	if err != nil {
		var notFoundErr *storage.ErrFileNotFound
		if errors.As(err, &notFoundErr) {
			return module.Metadata{}, &ModuleNotFoundError{
				ModuleAddr: moduleAddr,
				Cause:      err,
			}
		}
		return module.Metadata{}, fmt.Errorf("failed to read module file %s (%w)", path, err)
	}
	var mod module.Metadata
	if err := json.Unmarshal(fileContents, &mod); err != nil {
		return module.Metadata{}, fmt.Errorf("failed to parse module metadata file %s (%w)", path, err)
	}
	return mod, nil
}
