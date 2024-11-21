// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"fmt"
	"regexp"
	"runtime"
)

// ClientPullConfig contains the optional parameters to pull an OCI image.
type ClientPullConfig struct {
	// GOOS represents the operating system in Go terminology. The OCIClient uses this field to filter for the
	// correct image when encountering a multi-arch image. It defaults to empty, indicating that the current
	// runtime.GOOS should be used.
	GOOS string `json:"GOOS"`
	// GOARCH represents the CPU architecture in Go terminology. The OCIClient uses this field to filter for the
	// correct image when encountering a multi-arch image. It defaults to empty, indicating that the current
	// runtime.GOARCH should be used.
	GOARCH string `json:"GOARCH"`
}

// goStringsRe is a basic sanity check for GOOS and GOARCH values.
var goStringsRe = regexp.MustCompile(`^[a-zA-Z0-9.\-]+$`)

func (c *ClientPullConfig) ApplyDefaultsAndValidate() error {
	if c.GOOS == "" {
		c.GOOS = runtime.GOOS
	}
	if !goStringsRe.MatchString(c.GOOS) {
		return fmt.Errorf("invalid GOOS: %s (must match %s)", c.GOOS, goStringsRe.String())
	}
	if c.GOARCH == "" {
		c.GOARCH = runtime.GOARCH
	}
	if !goStringsRe.MatchString(c.GOARCH) {
		return fmt.Errorf("invalid GOARCH: %s (must match %s)", c.GOARCH, goStringsRe.String())
	}
	return nil
}

// ClientPullOpt configures an individual image pull. See OCIClient.PullImage
type ClientPullOpt func(c *ClientPullConfig) error

var goosRe = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
var goarchRe = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

// WithGOOS sets the operating system filter in Go terminology. The OCIClient uses this field to filter for the
// correct image when encountering a multi-arch image. It defaults to empty, indicating that the current
// runtime.GOOS should be used.
func WithGOOS(goos string) ClientPullOpt {
	return func(c *ClientPullConfig) error {
		if !goosRe.MatchString(goos) {
			return newPullConfigurationError(fmt.Sprintf("Invalid GOOS: %s (must match %s)", goos, goosRe), nil)
		}
		c.GOOS = goos
		return nil
	}
}

// WithGOARCH sets the CPU architecture filter in Go terminology. The OCIClient uses this field to filter for the
// correct image when encountering a multi-arch image. It defaults to empty, indicating that the current
// runtime.GOARCH should be used.
func WithGOARCH(goarch string) ClientPullOpt {
	return func(c *ClientPullConfig) error {
		if !goosRe.MatchString(goarch) {
			return newPullConfigurationError(fmt.Sprintf("Invalid GOARCH: %s (must match %s)", goarch, goarchRe), nil)
		}
		c.GOARCH = goarch
		return nil
	}
}
