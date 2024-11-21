// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// OCIRegistry describes a hostname and optionally a port for an OCI registry.
type OCIRegistry string

var hasPortRe = regexp.MustCompile(`:[0-9]+$`)

func (o OCIRegistry) Validate() error {
	if o == "" {
		return fmt.Errorf("empty OCI rgistry host")
	}
	if hasPortRe.MatchString(string(o)) {
		_, _, err := net.SplitHostPort(string(o))
		if err != nil {
			return fmt.Errorf("invalid OCI registry host: %s (%w)", o, err)
		}
	}
	_, _, err := net.SplitHostPort(string(o) + ":0")
	if err != nil {
		return fmt.Errorf("invalid OCI registry host: %s (%w)", o, err)
	}
	return nil
}

func (o OCIRegistry) Equals(other OCIRegistry) bool {
	return strings.EqualFold(string(o), string(other))
}

var _ validatable = OCIRegistry("")
