// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

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

func (r registryDataAPI) PutProviderNamespaceKey(ctx context.Context, namespace string, key provider.Key) error {
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

		k, e := crypto.NewKeyFromArmored(string(fileContents))
		if e != nil {
			return fmt.Errorf("failed to parse key file %s (%w)", p, e)
		}
		if strings.ToUpper(k.GetHexKeyID()) == key.KeyID {
			if e2 := r.storageAPI.PutFile(ctx, storage.Path(p), []byte(key.ASCIIArmor)); e2 != nil {
				return fmt.Errorf("failed to write key file %s (%w)", p, err)
			}
			return nil
		}
	}

	p := path.Join(basePath, key.KeyID+".asc")
	if e := r.storageAPI.PutFile(ctx, storage.Path(p), []byte(key.ASCIIArmor)); e != nil {
		return fmt.Errorf("failed to write key file %s (%w)", p, e)
	}
	return nil
}
