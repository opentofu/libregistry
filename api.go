// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package libregistry

import (
	"context"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/module"
	"github.com/opentofu/libregistry/vcs"
)

// API describes the API interface for accessing the registry.
type API interface {
	// AddModule adds a module based on a VCS repository. The VCS repository name must follow the naming convention
	// of the VCS implementation passed to the registry API on initialization.
	AddModule(ctx context.Context, vcsRepository string) error
	// UpdateModule updates the list of available versions for a module in the registry from its source repository.
	// This function is idempotent and adds the module to the storage if it does not exist yet.
	UpdateModule(ctx context.Context, moduleAddr module.Addr) error
}

// New creates a new instance of the registry API with the given GitHub client and data API instance.
func New(vcsClient vcs.Client, dataAPI metadata.ModuleDataAPI) (API, error) {
	return &api{
		dataAPI,
		vcsClient,
	}, nil
}

type api struct {
	dataAPI   metadata.ModuleDataAPI
	vcsClient vcs.Client
}
