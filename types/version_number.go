package types

type VersionNumber interface {
	// Parse the version number into major, minor, patch, stability, and stability number.
	Parse() (major int, minor int, patch int, stability string, stabilityNumber int, err error)
}
