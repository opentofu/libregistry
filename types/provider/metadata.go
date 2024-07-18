package provider

// Metadata contains information about the provider.
type Metadata struct {
	CustomRepository string    `json:"repository,omitempty"` // Optional. Custom repository from which to fetch the provider's metadata.
	Versions         []Version `json:"versions"`             // A list of version data, for each supported provider version.
}
