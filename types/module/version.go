package module

// Version represents a single version of a module.
type Version struct {
	// Version number of the provider. Correlates to a tag in the module repository.
	Version VersionNumber `json:"version"`
}

func (v Version) Normalize() Version {
	return Version{
		Version: v.Version.Normalize(),
	}
}