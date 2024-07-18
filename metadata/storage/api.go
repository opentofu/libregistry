// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
)

// API describes the necessary functions required to access the underlying files of the registry. You can
// implement this with a local filesystem, or a remote object storage.
type API interface {
	// ListFiles lists all file names in the specified directory. It must also list files created using PutFile that
	// have not yet been committed.
	ListFiles(ctx context.Context, directory Path) ([]string, error)
	// ListDirectories lists all subdirectories of the specified directory. It must also list all directories created
	// with PutFile before they are committed.
	ListDirectories(ctx context.Context, directory Path) ([]string, error)
	// PutFile creates a file at the specified path with the specified contents, making sure that all directories
	// for this file also exist. If the backing storage does not support directories, it must emulate directories.
	PutFile(ctx context.Context, path Path, contents []byte) error
	// GetFile returns the contents of the file at the specified path, or ErrFileNotFound if the file was not found,
	// or another error if there was a problem retrieving the file. The implementation must return files already
	// created by PutFile but not yet commited.
	GetFile(ctx context.Context, path Path) ([]byte, error)
	// DeleteFile removes the file from the backing storage on commit. If the file does not exist, it will not return
	// an error.
	DeleteFile(ctx context.Context, path Path) error
}
