// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package vcs

import (
	"regexp"
)

type Version string

var versionRe = regexp.MustCompile("^[a-zA-Z0-9/._-]+$")

const maxVersionLength = 255

// Validate validates the version strings against the assumptions the registry makes about version numbers.
func (v Version) Validate() error {
	if len(v) > maxVersionLength {
		return &InvalidVersionError{
			Version: v,
		}
	}
	if !versionRe.MatchString(string(v)) {
		return &InvalidVersionError{
			Version: v,
		}
	}
	return nil
}

type InvalidVersionError struct {
	Version Version
	Cause   error
}

func (r InvalidVersionError) Error() string {
	if r.Cause != nil {
		return "Failed to parse version: " + string(r.Version) + " (" + r.Cause.Error() + ")"
	}
	return "Failed to parse version: " + string(r.Version)
}

func (r InvalidVersionError) Unwrap() error {
	return r.Cause
}

type VersionNotFoundError struct {
	RepositoryAddr RepositoryAddr
	Version        Version
	Cause          error
}

func (v VersionNotFoundError) Error() string {
	if v.Cause != nil {
		return "Version " + string(v.Version) + " not found: in repository" + v.RepositoryAddr.String() + " (" + v.Cause.Error() + ")"
	}
	return "Version " + string(v.Version) + " not found: in repository" + v.RepositoryAddr.String()
}

func (v VersionNotFoundError) Unwrap() error {
	return v.Cause
}
