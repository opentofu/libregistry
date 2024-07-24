// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"context"

	"github.com/opentofu/libregistry/types/module"
)

// ModuleDataAPI provides access to the module data storage.
type ModuleDataAPI interface {
	// ListModules lists all modules in the registry.
	ListModules(ctx context.Context) ([]module.Addr, error)
	// ListModulesByNamespace returns a list of modules in a given namespace.
	ListModulesByNamespace(ctx context.Context, namespace string) ([]module.Addr, error)
	// ListModulesByNamespaceAndName returns a list of modules in the given namespace and name.
	ListModulesByNamespaceAndName(ctx context.Context, namespace string, name string) ([]module.Addr, error)

	// GetModule returns the current, uncommitted state of a module.
	GetModule(ctx context.Context, moduleAddr module.Addr) (module.Metadata, error)

	// GetAllModules returns a map of all module addresses and the metadata.
	GetAllModules(ctx context.Context) (map[module.Addr]module.Metadata, error)

	// PutModule queues the addition of a module with the given metadata. Call Commit() to write the changes
	// to the backing storage.
	PutModule(ctx context.Context, moduleAddr module.Addr, metadata module.Metadata) error
	// DeleteModule queues up the deletion of a given module.
	DeleteModule(ctx context.Context, moduleAddr module.Addr) error
}

const modulesDirectory = "modules"
