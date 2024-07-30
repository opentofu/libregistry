// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package vcs

import (
	"regexp"
	"time"
)

type VersionNumber string

var versionRe = regexp.MustCompile("^[a-zA-Z0-9/._-]+$")

const maxVersionLength = 255

// Validate validates the version strings against the assumptions the registry makes about version numbers.
func (v VersionNumber) Validate() error {
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

func (v VersionNumber) Equals(tag VersionNumber) bool {
	return v == tag
}

type Version struct {
	VersionNumber VersionNumber
	Created       time.Time
}

func (v Version) Validate() error {
	if err := v.VersionNumber.Validate(); err != nil {
		return err
	}
	return nil
}

func (v Version) Equals(other Version) bool {
	return v.VersionNumber.Equals(other.VersionNumber) && v.Created == other.Created
}

func (v Version) String() string {
	return string(v.VersionNumber)
}

type InvalidVersionError struct {
	Version VersionNumber
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
	Version        VersionNumber
	Cause          error
}

func (v VersionNotFoundError) Error() string {
	if v.Cause != nil {
		return "VersionNumber " + string(v.Version) + " not found: in repository" + v.RepositoryAddr.String() + " (" + v.Cause.Error() + ")"
	}
	return "VersionNumber " + string(v.Version) + " not found: in repository" + v.RepositoryAddr.String()
}

func (v VersionNotFoundError) Unwrap() error {
	return v.Cause
}
