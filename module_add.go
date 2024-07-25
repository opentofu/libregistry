// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package libregistry

import (
	"context"

	"github.com/opentofu/libregistry/types/module"
)

func (m api) AddModule(ctx context.Context, repository string) error {
	githubRepository, err := m.vcsClient.ParseRepositoryAddr(repository)
	if err != nil {
		return err
	}

	submitted, err := module.AddrFromRepository(githubRepository)
	if err != nil {
		return err
	}

	if err := submitted.Validate(); err != nil {
		return &ModuleAddFailedError{
			submitted,
			err,
		}
	}

	modules, err := m.dataAPI.ListModules(ctx)
	if err != nil {
		return err
	}

	for _, p := range modules {
		if p.Equals(submitted) {
			return &ModuleAlreadyExistsError{submitted}
		}
	}

	return m.UpdateModule(ctx, submitted)
}
