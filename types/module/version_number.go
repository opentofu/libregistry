package module

import (
	"strings"
)

// VersionNumber describes the semver version number.
type VersionNumber string

func (v VersionNumber) Normalize() VersionNumber {
	return VersionNumber("v" + strings.TrimPrefix(string(v), "v"))
}
