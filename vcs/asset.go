package vcs

import (
	"regexp"
)

type AssetName string

var assetNameRe = regexp.MustCompile("^[a-zA-Z0-9 ._-]+$")

const maxAssetNameLength = 255

// Validate validates asset names against assumptions the registry make about VCS assets.
func (a AssetName) Validate() error {
	if len(a) > maxAssetNameLength {
		return &InvalidAssetNameError{
			AssetName: a,
		}
	}
	if !assetNameRe.MatchString(string(a)) {
		return &InvalidAssetNameError{
			AssetName: a,
		}
	}
	return nil
}

type InvalidAssetNameError struct {
	AssetName AssetName
	Cause     error
}

func (r InvalidAssetNameError) Error() string {
	if r.Cause != nil {
		return "Failed to parse asset name: " + string(r.AssetName) + " (" + r.Cause.Error() + ")"
	}
	return "Failed to parse asset name: " + string(r.AssetName)
}

func (r InvalidAssetNameError) Unwrap() error {
	return r.Cause
}

type AssetNotFoundError struct {
	RepositoryAddr RepositoryAddr
	Version        Version
	Asset          AssetName
	Cause          error
}

func (a AssetNotFoundError) Error() string {
	if a.Cause != nil {
		return "Asset " + string(a.Asset) + " not found in version " + string(a.Version) + " of repository" + a.RepositoryAddr.String() + " (" + a.Cause.Error() + ")"
	}
	return "Asset " + string(a.Asset) + " not found in version " + string(a.Version) + " of repository" + a.RepositoryAddr.String()
}

func (a AssetNotFoundError) Unwrap() error {
	return a.Cause
}
