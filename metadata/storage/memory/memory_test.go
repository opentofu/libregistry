// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package memory_test

import (
	"testing"

	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/metadata/storage/memory"
)

func TestFileHandling(t *testing.T) {
	storage.TestStorageAPI(t, func(t *testing.T) storage.API {
		return memory.New()
	})
}
