// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package vcs

import (
	"regexp"
)

// nameRe is a common understanding of acceptable repository addresses. This may need to be changed later if it turns
// out that other VCS' support different address styles.
var nameRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// RepositoryAddr holds a reference to a repository. For simplicity, the current system does not support more complex
// URL structures.
type RepositoryAddr struct {
	Org OrganizationAddr
	// Name is the URL fragment of a repository.
	Name string
}

func (r RepositoryAddr) String() string {
	return string(r.Org) + "/" + r.Name
}

// Validate checks the assumptions the registry makes about repositories.
func (r RepositoryAddr) Validate() error {
	if err := r.Org.Validate(); err != nil {
		return &InvalidRepositoryAddrError{RepositoryAddr: r, Cause: err}
	}
	if !nameRe.MatchString(r.Name) {
		return &InvalidRepositoryAddrError{
			RepositoryAddr: r,
		}
	}
	return nil
}

type InvalidRepositoryAddrError struct {
	RepositoryString string
	RepositoryAddr   RepositoryAddr
	Cause            error
}

func (r InvalidRepositoryAddrError) Error() string {
	if r.Cause != nil {
		if r.RepositoryString != "" {
			return "Failed to parse repository address: " + string(r.RepositoryString) + " (" + r.Cause.Error() + ")"
		}
		return "Failed to parse repository address: " + string(r.RepositoryAddr.String()) + " (" + r.Cause.Error() + ")"
	}
	if r.RepositoryString != "" {
		return "Failed to parse repository address: " + string(r.RepositoryString)
	}
	return "Failed to parse repository address: " + string(r.RepositoryAddr.String())
}

func (r InvalidRepositoryAddrError) Unwrap() error {
	return r.Cause
}

type RepositoryNotFoundError struct {
	RepositoryAddr RepositoryAddr
	Cause          error
}

func (r RepositoryNotFoundError) Error() string {
	if r.Cause != nil {
		return "Repository not found: " + r.RepositoryAddr.String() + " (" + r.Cause.Error() + ")"
	}
	return "Repository not found: " + r.RepositoryAddr.String()
}

func (r RepositoryNotFoundError) Unwrap() error {
	return r.Cause
}
