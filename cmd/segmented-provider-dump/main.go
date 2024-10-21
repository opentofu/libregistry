// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

// Package main contains a tool to dump providers according to version batch sizes.
package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/filesystem"
)

func main() {
	if len(os.Args) != 4 {
		_, _ = os.Stderr.Write([]byte("Usage: segmented-provider-dump path/to/registry batchsize destination"))
		os.Exit(1)
	}

	meta, err := metadata.New(filesystem.New(os.Args[1]))
	if err != nil {
		_, _ = os.Stderr.Write([]byte(fmt.Errorf("failed to initialize metadata system; did you pass the correct registry directory? (%w)", err).Error()))
		os.Exit(1)
	}
	batchSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}
	destination := os.Args[3]
	stat, err := os.Stat(destination)
	if err != nil {
		_, _ = os.Stderr.Write([]byte("Invalid destination: " + destination + " (" + err.Error() + ")"))
		os.Exit(1)
	}
	if !stat.IsDir() {
		_, _ = os.Stderr.Write([]byte("Invalid destination: " + destination + " (not a directory)"))
		os.Exit(1)
	}

	ctx := context.Background()

	providers, err := meta.ListProviders(ctx, true)
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}

	versions := 0
	batch := 0
	var fh *os.File

	openBatch := func() error {
		fh, err = os.Create(path.Join(destination, fmt.Sprintf("%d.txt", batch)))
		if err != nil {
			return fmt.Errorf("failed to open batch file %d.txt", batch)
		}
		return nil
	}
	closeBatch := func() error {
		if fh != nil {
			if err := fh.Close(); err != nil {
				return fmt.Errorf("failed to close batch file %d.txt", batch)
			}
		}
		return nil
	}
	if err := openBatch(); err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}
	for _, providerAddr := range providers {
		provider, err := meta.GetProvider(ctx, providerAddr, true)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(err.Error()))
			os.Exit(1)
		}
		if versions+len(provider.Versions) > batchSize {
			if err := closeBatch(); err != nil {
				_, _ = os.Stderr.Write([]byte(fmt.Sprintf("Failed to close %d.txt (%v)", batch, err)))
				os.Exit(1)
			}
			versions = 0
			batch += 1
			if err := openBatch(); err != nil {
				_, _ = os.Stderr.Write([]byte(fmt.Sprintf("Failed to open %d.txt (%v)", batch, err)))
				os.Exit(1)
			}
		}
		versions += len(provider.Versions)
		if _, err := fh.Write([]byte(fmt.Sprintf("%s\n", providerAddr))); err != nil {
			_, _ = os.Stderr.Write([]byte(fmt.Sprintf("Failed to write %d.txt (%v)", batch, err)))
			_ = closeBatch()
			os.Exit(1)
		}
	}
	if err := closeBatch(); err != nil {
		_, _ = os.Stderr.Write([]byte(fmt.Sprintf("Failed to close %d.txt (%v)", batch, err)))
		os.Exit(1)
	}
}
