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

func (r registryDataAPI) GetProviderNamespaceKey(ctx context.Context, namespace string, keyID string) (provider.Key, error) {
	namespace = provider.NormalizeNamespace(namespace)
	basePath := path.Join(keysDirectory, namespace[0:1], namespace)
	files, err := r.storageAPI.ListFiles(ctx, storage.Path(basePath))
	if err != nil {
		return provider.Key{}, fmt.Errorf("failed to list keys for namespace %s (%w)", namespace, err)
	}

	for _, file := range files {
		p := path.Join(basePath, file)
		fileContents, e := r.storageAPI.GetFile(ctx, storage.Path(p))
		if e != nil {
			return provider.Key{}, fmt.Errorf("failed to read file %s (%w)", p, e)
		}

		armored := string(fileContents)

		key, e := crypto.NewKeyFromArmored(armored)
		if e != nil {
			return provider.Key{}, fmt.Errorf("failed to parse key file %s (%w)", p, e)
		}
		if strings.ToUpper(key.GetHexKeyID()) == keyID {
			return provider.Key{
				ASCIIArmor: armored,
				KeyID:      keyID,
			}, nil
		}
	}

	// Key not found.
	return provider.Key{}, fmt.Errorf("cannot find key ID %s for provider namespace %s", keyID, namespace)
}
