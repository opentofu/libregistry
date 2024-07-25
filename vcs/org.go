// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package vcs

import (
	"regexp"
)

// OrganizationAddr refers to an organization within a VCS system.
type OrganizationAddr string

// Validate validates the organization name against
func (o OrganizationAddr) Validate() error {
	if len(o) > maxOrgNameLength {
		return &InvalidOrganizationAddrError{
			OrganizationAddr: o,
		}
	}
	if !orgNameRe.MatchString(string(o)) {
		return &InvalidOrganizationAddrError{
			OrganizationAddr: o,
		}
	}
	return nil
}

func (o OrganizationAddr) String() string {
	return string(o)
}

type InvalidOrganizationAddrError struct {
	OrganizationAddr OrganizationAddr
	Cause            error
}

func (r InvalidOrganizationAddrError) Error() string {
	if r.Cause != nil {
		return "Invalid organization name: " + string(r.OrganizationAddr) + " (" + r.Cause.Error() + ")"
	}
	return "Invalid organization name: " + string(r.OrganizationAddr)
}

func (r InvalidOrganizationAddrError) Unwrap() error {
	return r.Cause
}

type OrganizationNotFoundError struct {
	OrganizationAddr OrganizationAddr
	Cause            error
}

func (r OrganizationNotFoundError) Error() string {
	if r.Cause != nil {
		return "Organization not found: " + r.OrganizationAddr.String() + " (" + r.Cause.Error() + ")"
	}
	return "Organization not found: " + r.OrganizationAddr.String()
}

func (r OrganizationNotFoundError) Unwrap() error {
	return r.Cause
}

var orgNameRe = regexp.MustCompile("^([a-zA-Z0-9]+)(|(-[a-zA-Z0-9]+)*)$")

const maxOrgNameLength = 255
