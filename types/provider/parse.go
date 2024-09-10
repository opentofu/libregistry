// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"

	"github.com/opentofu/libregistry/vcs"
)

var providerRepoRe = regexp.MustCompile("terraform-provider-(?P<Name>[a-zA-Z0-9-]*)")

// AddrFromRepository parses a repository name and returns a provider address from it if valid.
func AddrFromRepository(repository vcs.RepositoryAddr) (Addr, error) {
	match := providerRepoRe.FindStringSubmatch(repository.Name)
	if match == nil {
		return Addr{}, fmt.Errorf("invalid provider repository name: %s", repository.String())
	}

	return Addr{
		Namespace: string(repository.Org),
		Name:      match[providerRepoRe.SubexpIndex("Name")],
	}, nil
}

// VersionFromVCS converts a vcs.VersionNumber into a VersionNumber.
func VersionFromVCS(vcsVersion vcs.VersionNumber) (VersionNumber, error) {
	ver := VersionNumber(vcsVersion)
	return ver, ver.Validate()
}
