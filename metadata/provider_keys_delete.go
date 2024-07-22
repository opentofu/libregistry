// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) DeleteProviderNamespaceKey(ctx context.Context, namespace string, keyID string) error {
	namespace = provider.NormalizeNamespace(namespace)
	basePath := path.Join(keysDirectory, namespace[0:1], namespace)
	files, err := r.storageAPI.ListFiles(ctx, storage.Path(basePath))
	if err != nil {
		return fmt.Errorf("failed to list keys for namespace %s (%w)", namespace, err)
	}

	for _, file := range files {
		p := path.Join(basePath, file)
		fileContents, e := r.storageAPI.GetFile(ctx, storage.Path(p))
		if e != nil {
			return fmt.Errorf("failed to read file %s (%w)", p, e)
		}

		key, e := crypto.NewKeyFromArmored(string(fileContents))
		if e != nil {
			return fmt.Errorf("failed to parse key file %s (%w)", p, e)
		}
		if strings.ToUpper(key.GetHexKeyID()) == keyID {
			if e2 := r.storageAPI.DeleteFile(ctx, storage.Path(p)); e2 != nil {
				return fmt.Errorf("failed to delete key file %s (%w)", p, e2)
			} else {
				return nil
			}
		}
	}

	// Key not found.
	return nil
}
