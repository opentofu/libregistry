// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package filesystem_test

import (
	"testing"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/metadata/storage/filesystem"
)

func TestFileHandling(t *testing.T) {
	storage.TestStorageAPI(t, func(t *testing.T) storage.API {
		return filesystem.New(t.TempDir())
	})
}
