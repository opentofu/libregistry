// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"fmt"
	"regexp"

	"github.com/opentofu/libregistry/vcs"
)

var moduleRepoRe = regexp.MustCompile("terraform-(?P<Target>[a-zA-Z0-9]*)-(?P<Name>[a-zA-Z0-9-]*)")

// AddrFromRepository parses a repository name and returns a module address from it if valid.
func AddrFromRepository(repository vcs.RepositoryAddr) (Addr, error) {
	match := moduleRepoRe.FindStringSubmatch(repository.Name)
	if match == nil {
		return Addr{}, fmt.Errorf("invalid module repository name: %s", repository.String())
	}

	return Addr{
		Namespace:    string(repository.Org),
		Name:         match[moduleRepoRe.SubexpIndex("Name")],
		TargetSystem: match[moduleRepoRe.SubexpIndex("Target")],
	}, nil
}

// VersionFromVCS converts a vcs.Version into a VersionNumber.
func VersionFromVCS(vcsVersion vcs.Version) (VersionNumber, error) {
	ver := VersionNumber(vcsVersion)
	return ver, ver.Validate()
}
