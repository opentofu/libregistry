package provider

// Version contains information about a specific provider version.
type Version struct {
	Version             string   `json:"version"`               // The version number of the provider.
	Protocols           []string `json:"protocols"`             // The protocol versions the provider supports.
	SHASumsURL          string   `json:"shasums_url"`           // The URL to the SHA checksums file.
	SHASumsSignatureURL string   `json:"shasums_signature_url"` // The URL to the GPG signature of the SHA checksums file.
	Targets             []Target `json:"targets"`               // A list of target platforms for which this provider version is available.
}
