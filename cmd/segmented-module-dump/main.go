// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

// Package main contains a tool to dump modules according to version batch sizes.
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/filesystem"
)

func main() {
	if len(os.Args) != 3 {
		_, _ = os.Stderr.Write([]byte("Usage: segmented-module-dump path/to/registry batchsize"))
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

	ctx := context.Background()

	providers, err := meta.ListModules(ctx)
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}

	versions := 0
	batch := 0
	fmt.Printf("Provider\tVersions\tBatch\n")
	for _, moduleAddr := range providers {
		module, err := meta.GetModule(ctx, moduleAddr)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(err.Error()))
			os.Exit(1)
		}
		if versions+len(module.Versions) > batchSize {
			versions = 0
			batch += 1
		}
		versions += len(module.Versions)
		fmt.Printf("%s\t%d\t%d\n", moduleAddr, len(module.Versions), batch)
	}
}
