package provider_key_verifier

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/opentofu/libregistry/logger"
	"github.com/opentofu/libregistry/types/provider"
)

// Config holds the configuration for GitHub.
type Config struct {
	// Logger holds the logger to write any logs to.
	Logger logger.Logger
	// HTTPClient holds the HTTP client to use for API requests. Note that this only affects API and RSS feed requests,
	// but not git clone commands as those are done using the command line.
	HTTPClient *http.Client
	// Number of versions that are going to be checked if they were signed
	NumVersionsToCheck uint8
	checkFn            CheckFn
}

// Opt is a function that modifies the config.
type Opt func(config *Config) error
type CheckFn func(pkv providerKeyVerifier, ctx context.Context, version provider.Version) error

// ApplyDefaults adds the default values if none are present.
func (c *Config) ApplyDefaults() {
	if c.Logger == nil {
		c.Logger = logger.NewNoopLogger()
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
		transport := http.DefaultTransport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
		c.HTTPClient.Transport = transport
	}
}

// WithVersionsToCheck is a functional option to set the number of versions to check for a provider.
func WithNumVersionsToCheck(versionsToCheck uint8) Opt {
	return func(config *Config) error {
		config.NumVersionsToCheck = versionsToCheck
		return nil
	}
}

// WithLogger is a functional option to set the logger
func WithLogger(logger logger.Logger) Opt {
	return func(config *Config) error {
		config.Logger = logger
		return nil
	}
}

// WithHTTPClient is a functional option to set the http Client
func WithHTTPClient(httpClient *http.Client) Opt {
	return func(config *Config) error {
		config.HTTPClient = httpClient
		return nil
	}
}

// WithCheckFn is a functional option to set the function used to check the provider version
func WithCheckFn(checkFn CheckFn) Opt {
	return func(config *Config) error {
		config.checkFn = checkFn
		return nil
	}
}
