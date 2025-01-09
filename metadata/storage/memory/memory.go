// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package memory

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/opentofu/libregistry/metadata/storage"
)

// New returns an in-memory filesystem for testing purposes.
func New() storage.API {
	return &api{
		root: &directory{
			directories: map[string]*directory{},
			files:       map[string][]byte{},
		},
	}
}

type api struct {
	root *directory
}

type directory struct {
	directories map[string]*directory
	files       map[string][]byte
}

func (a *api) ListFiles(_ context.Context, directory storage.Path) ([]string, error) {
	if err := directory.Validate(); err != nil {
		return nil, err
	}
	current, err := a.resolveDirectory(directory)
	if err != nil {
		var e *storage.ErrFileNotFound
		if errors.As(err, &e) {
			return nil, nil
		}
		return nil, err
	}
	result := make([]string, len(current.files))
	i := 0
	for file := range current.files {
		result[i] = file
		i++
	}
	return result, nil
}

func (a *api) resolveDirectory(directory storage.Path) (*directory, error) {
	current := a.root
	if directory != "" {
		parts := strings.Split(strings.Trim(string(directory), "/"), "/")
		for _, part := range parts {
			subdir, ok := current.directories[part]
			if !ok {
				return nil, &storage.ErrFileNotFound{
					Path: directory,
				}
			}
			current = subdir
		}
	}
	return current, nil
}

func (a *api) ensureDirectory(path storage.Path) *directory {
	current := a.root
	if path != "" {
		parts := strings.Split(string(path), "/")
		for _, part := range parts {
			_, ok := current.directories[part]
			if !ok {
				current.directories[part] = &directory{
					directories: map[string]*directory{},
					files:       map[string][]byte{},
				}
			}
			current = current.directories[part]
		}
	}
	return current
}

func (a *api) ListDirectories(_ context.Context, directory storage.Path) ([]string, error) {
	if err := directory.Validate(); err != nil {
		return nil, err
	}

	current, err := a.resolveDirectory(directory)
	if err != nil {
		var notFound *storage.ErrFileNotFound
		if errors.As(err, &notFound) {
			return nil, nil
		}
		return nil, err
	}
	result := make([]string, len(current.directories))
	i := 0
	for file := range current.directories {
		result[i] = file
		i++
	}
	return result, nil
}

func (a *api) PutFile(_ context.Context, filePath storage.Path, contents []byte) error {
	if err := filePath.Validate(); err != nil {
		return err
	}

	current := a.ensureDirectory(filePath.Basename())

	if _, ok := current.directories[filePath.Filename()]; ok {
		return &storage.ErrFileAlreadyExists{
			Path: filePath,
		}
	}

	current.files[filePath.Filename()] = contents
	return nil
}

func (a *api) GetFile(_ context.Context, filePath storage.Path) ([]byte, error) {
	if err := filePath.Validate(); err != nil {
		return nil, err
	}

	current, err := a.resolveDirectory(filePath.Basename())
	if err != nil {
		return nil, err
	}
	contents, ok := current.files[filePath.Filename()]
	if !ok {
		return nil, &storage.ErrFileNotFound{
			Path: filePath,
		}
	}
	return contents, nil
}

func (a *api) FileExists(_ context.Context, filePath storage.Path) (bool, error) {
	if err := filePath.Validate(); err != nil {
		return false, err
	}

	current, err := a.resolveDirectory(filePath.Basename())
	if err != nil {
		var notFound *storage.ErrFileNotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	_, ok := current.files[filePath.Filename()]
	return ok, nil
}

func (a *api) DeleteFile(_ context.Context, filePath storage.Path) error {
	if err := filePath.Validate(); err != nil {
		return err
	}

	current, err := a.resolveDirectory(filePath.Basename())
	if err != nil {
		var e *storage.ErrFileNotFound
		if errors.As(err, &e) {
			return nil
		}
		return err
	}
	delete(current.files, filePath.Filename())
	return nil
}

func (a *api) DownloadFile(ctx context.Context, url string, filePath storage.Path) error {
	if err := filePath.Validate(); err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	current := a.ensureDirectory(filePath.Basename())

	if _, ok := current.directories[filePath.Filename()]; ok {
		return &storage.ErrFileAlreadyExists{
			Path: filePath,
		}
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	current.files[filePath.Filename()] = contents
	return nil
}
