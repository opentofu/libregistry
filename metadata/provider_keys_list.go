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
)

func (r registryDataAPI) ListProviderNamespacesWithKeys(_ context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r registryDataAPI) ListProviderNamespaceKeyIDs(ctx context.Context, namespace string) ([]string, error) {
	basePath := path.Join(keysDirectory, namespace[0:1], namespace)
	files, err := r.storageAPI.ListFiles(ctx, storage.Path(basePath))
	if err != nil {
		return nil, fmt.Errorf("failed to list keys for namespace %s (%w)", namespace, err)
	}

	var results []string

	for _, file := range files {
		p := path.Join(basePath, file)
		fileContents, e := r.storageAPI.GetFile(ctx, storage.Path(p))
		if e != nil {
			return nil, fmt.Errorf("failed to read file %s (%w)", p, e)
		}

		armored := string(fileContents)

		key, e := crypto.NewKeyFromArmored(armored)
		if e != nil {
			return nil, fmt.Errorf("failed to parse key file %s (%w)", p, e)
		}
		results = append(results, strings.ToUpper(key.GetHexKeyID()))
	}

	return results, nil
}
