// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/types/module"
)

func (r registryDataAPI) ListModules(ctx context.Context) ([]module.Addr, error) {
	moduleLetters, err := r.storageAPI.ListDirectories(ctx, "modules")
	if err != nil {
		// The modules directory does not exist:
		var notFoundError storage.ErrFileNotFound
		if errors.As(err, &notFoundError) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list 'modules' directory (%w)", err)
	}
	var results []module.Addr
	for _, letter := range moduleLetters {
		letterResults, e := r.listModulesByNamespaceLetter(ctx, letter)
		if e != nil {
			return nil, e
		}
		results = append(results, letterResults...)
	}
	return results, nil
}

func (r registryDataAPI) listModulesByNamespaceLetter(ctx context.Context, letter string) ([]module.Addr, error) {
	p := storage.Path(path.Join(modulesDirectory, letter))
	namespaces, e := r.storageAPI.ListDirectories(nil, p)
	if e != nil {
		// The letter directory does not exist:
		var notFoundError storage.ErrFileNotFound
		if errors.As(e, &notFoundError) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list module directory %s (%w)", p, e)
	}
	var results []module.Addr
	for _, namespace := range namespaces {
		namespaceResults, e2 := r.ListModulesByNamespace(ctx, namespace)
		if e2 != nil {
			return nil, e2
		}
		results = append(results, namespaceResults...)
	}
	return results, nil
}

func (r registryDataAPI) ListModulesByNamespace(ctx context.Context, namespace string) ([]module.Addr, error) {
	namespace = module.NormalizeNamespace(namespace)
	p := storage.Path(path.Join(modulesDirectory, namespace[0:1], namespace))
	directories, err := r.storageAPI.ListDirectories(nil, p)
	if err != nil {
		// The namespace directory does not exist:
		var notFoundError storage.ErrFileNotFound
		if errors.As(err, &notFoundError) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list namespace directory %s (%w)", p, err)
	}
	var results []module.Addr
	for _, name := range directories {
		mods, e := r.ListModulesByNamespaceAndName(ctx, namespace, name)
		if e != nil {
			return nil, e
		}
		results = append(results, mods...)
	}
	return results, nil
}

func (r registryDataAPI) ListModulesByNamespaceAndName(ctx context.Context, namespace string, name string) ([]module.Addr, error) {
	namespace = module.NormalizeNamespace(namespace)
	name = module.NormalizeName(name)
	p := storage.Path(path.Join(modulesDirectory, namespace[0:1], namespace, name))
	files, err := r.storageAPI.ListFiles(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in module name directory %s (%w)", p, err)
	}
	var result []module.Addr
	for _, file := range files {
		if strings.HasSuffix(file, ".json") {
			result = append(result, module.Addr{
				Namespace:    namespace,
				Name:         name,
				TargetSystem: strings.ToLower(strings.TrimSuffix(file, ".json")),
			})
		}
	}
	return result, nil
}
