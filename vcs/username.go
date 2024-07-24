// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package vcs

import (
	"regexp"
)

// Username is a user in a VCS system.
type Username string

// Validate validates the username against the assumptions the registry makes about usernames in VCS systems for safety
// when including them in URLs.
func (u Username) Validate() error {
	if len(u) > maxUsernameLength {
		return &InvalidUsernameError{Username: u}
	}
	if !usernameRe.MatchString(string(u)) {
		return &InvalidUsernameError{Username: u}
	}
	return nil
}

// InvalidUsernameError indicates that a given username did not meet the assumptions we make about VCS usernames. This
// may be due to an incorrect assumption or legitimately invalid data.
type InvalidUsernameError struct {
	Username Username
	Cause    error
}

func (i InvalidUsernameError) Error() string {
	if i.Username == "" {
		if i.Cause != nil {
			return "Empty username (" + i.Cause.Error() + ")"
		}
		return "Empty username"
	}
	if i.Cause != nil {
		return "Invalid username: " + string(i.Username) + " (" + i.Cause.Error() + ")"
	}
	return "Invalid username: " + string(i.Username)
}

var usernameRe = regexp.MustCompile("^([a-zA-Z0-9]+)(|(-[a-zA-Z0-9]+)*)$")

const maxUsernameLength = 255
