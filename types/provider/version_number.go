// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/opentofu/libregistry/vcs"
	"golang.org/x/mod/semver"
)

const majorName = "major"
const minorName = "minor"
const patchName = "patch"
const stabilityName = "stability"
const stabilityNumberName = "stabilityNumber"

var versionRe = regexp.MustCompile(`^v(?P<` + majorName + `>[0-9]+)\.(?P<` + minorName + `>[0-9]+)\.(?P<` + patchName + `>[0-9]+)(|-(?P<` + stabilityName + `>[a-zA-Z.-]+)(?P<` + stabilityNumberName + `>[0-9]+))$`)
var versionReNames = map[string]int{}

func init() {
	for index, name := range versionRe.SubexpNames() {
		if name != "" {
			versionReNames[name] = index
		}
	}
}

const maxVersionLength = 255

// VersionNumber describes the semver version number.
//
// swagger:model ProviderVersionNumber
type VersionNumber string

func (v VersionNumber) Normalize() VersionNumber {
	return VersionNumber("v" + strings.TrimPrefix(string(v), "v"))
}

func (v VersionNumber) Compare(other VersionNumber) int {
	return semver.Compare(string(v.Normalize()), string(other.Normalize()))
}

func (v VersionNumber) Validate() error {
	normalizedV := v.Normalize()
	if len(normalizedV) > maxVersionLength {
		return &InvalidVersionNumber{v}
	}
	if !versionRe.MatchString(string(normalizedV)) {
		return &InvalidVersionNumber{v}
	}
	return nil
}

// ToVCSVersion creates a vcs.VersionNumber from the VersionNumber. Call ToVCSVersion() before you call Normalize() in
// order to get the correct VCS version.
func (v VersionNumber) ToVCSVersion() vcs.VersionNumber {
	return vcs.VersionNumber(v)
}

func (v VersionNumber) Parse() (major int, minor int, patch int, stability string, stabilityNumber int, err error) {
	versionReNames := versionReNames
	submatches := versionRe.FindStringSubmatch(string(v.Normalize()))
	if len(submatches) == 0 {
		return 0, 0, 0, "", 0, fmt.Errorf("failed to parse version (must match %s)", versionRe.String())
	}
	numbers := map[string]int{}
	for _, matchName := range []string{majorName, minorName, patchName, stabilityNumberName} {
		value := submatches[versionReNames[matchName]]
		if matchName == stabilityNumberName && value == "" {
			numbers[matchName] = 0
			continue
		}
		numbers[matchName], err = strconv.Atoi(value)
		if err != nil {
			return 0, 0, 0, "", 0, fmt.Errorf("failed to parse version (%s is invalid: %w)", matchName, err)
		}
	}
	return numbers[majorName], numbers[minorName], numbers[patchName], submatches[versionReNames[stabilityName]], numbers[stabilityNumberName], nil
}

type InvalidVersionNumber struct {
	VersionNumber VersionNumber
}

func (i InvalidVersionNumber) Error() string {
	return "Invalid version: " + string(i.VersionNumber)
}
