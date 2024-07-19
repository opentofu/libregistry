// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"testing"
)

// TestStorageAPI provides a test tool for storage implementations.
func TestStorageAPI(t *testing.T, factory func(t *testing.T) API) {
	const testFile1 = "test.txt"
	const testFile2Directory = "test"
	const testFile2Name = "test.txt"
	const testFile2 = testFile2Directory + "/" + testFile2Name
	var testFileContents = []byte("Hello world!")

	ctx := context.Background()

	t.Run("root", func(t *testing.T) {
		fa := factory(t)
		files, err := fa.ListFiles(ctx, "")
		if err != nil {
			t.Fatalf("Cannot list root directory files (%v)", err)
		}
		if len(files) != 0 {
			t.Fatalf("Root directory is not empty on start (%d items)", len(files))
		}

		if err := fa.PutFile(ctx, testFile1, testFileContents); err != nil {
			t.Fatalf("Cannot put file %s (%v)", testFile1, err)
		}

		files, err = fa.ListFiles(ctx, "")
		if err != nil {
			t.Fatalf("Cannot list root directory files (%v)", err)
		}
		if len(files) != 1 {
			t.Fatalf("Unexpected file count in root directory (%d)", len(files))
		}
		if files[0] != testFile1 {
			t.Fatalf("Unexpected file name: %s (want: %s)", files[0], testFile1)
		}

		contents, err := fa.GetFile(ctx, testFile1)
		if err != nil {
			t.Fatalf("Failed to fetch test file %s (%v)", testFile1, err)
		}
		if string(contents) != string(testFileContents) {
			t.Fatalf("Incorrect file contents: %s (%v)", contents, err)
		}

		if err := fa.DeleteFile(ctx, testFile1); err != nil {
			t.Fatalf("Failed to delete test file (%v)", err)
		}

		if err := fa.DeleteFile(ctx, testFile1); err != nil {
			t.Fatalf("Failed to delete already-deleted file (%v)", err)
		}

		if _, err = fa.GetFile(ctx, testFile1); err == nil {
			t.Fatalf("Fetched already-deleted file (%s)", testFile1)
		}
	})

	t.Run("subdir", func(t *testing.T) {
		fa := factory(t)
		files, err := fa.ListFiles(ctx, testFile2Directory)
		if err != nil {
			t.Fatalf("Cannot list %s directory files (%v)", testFile2Directory, err)
		}
		if len(files) != 0 {
			t.Fatalf("Root directory is not empty on start (%d items)", len(files))
		}

		if err := fa.PutFile(ctx, testFile2, testFileContents); err != nil {
			t.Fatalf("Cannot put file %s (%v)", testFile2, err)
		}

		files, err = fa.ListFiles(ctx, testFile2Directory)
		if err != nil {
			t.Fatalf("Cannot list the %s directory files (%v)", testFile2Directory, err)
		}
		if len(files) != 1 {
			t.Fatalf("Unexpected file count in the %s directory (%d)", testFile2Directory, len(files))
		}
		if files[0] != testFile2Name {
			t.Fatalf("Unexpected file name: %s (want: %s)", files[0], testFile2)
		}

		files, err = fa.ListFiles(ctx, "")
		if err != nil {
			t.Fatalf("Cannot list root directory files (%v)", err)
		}
		if len(files) != 0 {
			t.Fatalf("Unexpected file count in the root directory (%d)", len(files))
		}

		directories, err := fa.ListDirectories(ctx, "")
		if err != nil {
			t.Fatalf("Cannot list root directory subdirectories (%v)", err)
		}
		if len(directories) != 1 {
			t.Fatalf("Unexpected directory count in the root directory (%d)", len(directories))
		}

		if directories[0] != testFile2Directory {
			t.Fatalf("unexpected directory name: %s", directories[0])
		}
		contents, err := fa.GetFile(ctx, testFile2)
		if err != nil {
			t.Fatalf("Failed to fetch test file %s (%v)", testFile2, err)
		}
		if string(contents) != string(testFileContents) {
			t.Fatalf("Incorrect file contents: %s (%v)", contents, err)
		}

		if err := fa.DeleteFile(ctx, testFile2); err != nil {
			t.Fatalf("Failed to delete test file (%v)", err)
		}

		if err := fa.DeleteFile(ctx, testFile2); err != nil {
			t.Fatalf("Failed to delete already-deleted file (%v)", err)
		}

		if _, err = fa.GetFile(ctx, testFile2); err == nil {
			t.Fatalf("Fetched already-deleted file (%s)", testFile2)
		}
	})
}
