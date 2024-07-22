// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package filesystem

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/opentofu/libregistry/metadata/storage"
)

// New creates an API implementation that works with a local filesystem.
func New(directory string) storage.API {
	return &storageAPI{
		directory: directory,
	}
}

type storageAPI struct {
	directory string
}

func (f *storageAPI) ListFiles(_ context.Context, directory storage.Path) ([]string, error) {
	if err := directory.Validate(); err != nil {
		return nil, err
	}

	dir, err := os.ReadDir(path.Join(f.directory, string(directory)))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
	}
	var result []string
	for _, entry := range dir {
		if !entry.IsDir() {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}

func (f *storageAPI) ListDirectories(_ context.Context, directory storage.Path) ([]string, error) {
	if err := directory.Validate(); err != nil {
		return nil, err
	}

	dir, err := os.ReadDir(path.Join(f.directory, string(directory.Basename())))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
	}
	var result []string
	for _, entry := range dir {
		if entry.IsDir() {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}

func (f *storageAPI) PutFile(_ context.Context, filePath storage.Path, contents []byte) error {
	if err := filePath.Validate(); err != nil {
		return err
	}

	fullPath := path.Join(f.directory, string(filePath.Basename()))
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create base directory %s (%w)", fullPath, err)
	}

	fullFilePath := path.Join(f.directory, string(filePath))
	if err := os.WriteFile(fullFilePath, contents, 0644); err != nil {
		if os.IsExist(err) {
			return storage.ErrFileAlreadyExists{
				Path: filePath,
			}
		}
		return fmt.Errorf("failed to create file %s (%w)", fullFilePath, err)
	}
	return nil
}

func (f *storageAPI) GetFile(_ context.Context, filePath storage.Path) ([]byte, error) {
	if err := filePath.Validate(); err != nil {
		return nil, err
	}

	fullPath := path.Join(f.directory, string(filePath))
	contents, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.ErrFileNotFound{
				Path: filePath,
			}
		} else {
			return nil, fmt.Errorf("cannot get file %s (%w)", fullPath, err)
		}
	}
	return contents, nil
}

func (f *storageAPI) FileExists(_ context.Context, filePath storage.Path) (bool, error) {
	if err := filePath.Validate(); err != nil {
		return false, err
	}

	fullPath := path.Join(f.directory, string(filePath))
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (f *storageAPI) DeleteFile(_ context.Context, filePath storage.Path) error {
	if err := filePath.Validate(); err != nil {
		return err
	}

	if filePath == "" {
		return storage.ErrFileNotFound{
			Path: filePath,
		}
	}

	fullPath := path.Join(f.directory, string(filePath))
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("cannot remove %s (%w)", fullPath, err)
	}
	return nil
}
