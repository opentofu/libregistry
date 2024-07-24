// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package libregistry

import (
	"context"
	"fmt"
	"regexp"

	"github.com/opentofu/libregistry/types/module"
)

var moduleRepoRe = regexp.MustCompile("terraform-(?P<Target>[a-zA-Z0-9]*)-(?P<Name>[a-zA-Z0-9-]*)")

func (m api) AddModule(ctx context.Context, repository string) error {
	githubRepository, err := m.vcsClient.ParseRepositoryAddr(repository)
	if err != nil {
		return err
	}

	match := moduleRepoRe.FindStringSubmatch(githubRepository.Name)
	if match == nil {
		return fmt.Errorf("invalid repository name: %s", githubRepository.String())
	}

	submitted := module.Addr{
		Namespace:    string(githubRepository.Org),
		Name:         match[moduleRepoRe.SubexpIndex("Name")],
		TargetSystem: match[moduleRepoRe.SubexpIndex("Target")],
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
