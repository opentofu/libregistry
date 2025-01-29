package provider_key_verifier

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

// ProviderKeyVerifier describes the functions for verifying if a key was used to sign a list of providers.
type ProviderKeyVerifier interface {
	// VerifyProvider verifies if the key was used to sign a provider addr. It returns a list of the valid versions signed by this key.
	VerifyProvider(ctx context.Context, provider provider.Addr) ([]*provider.Version, error)
}

// New creates a new instance of the provider key verification package with the given f keyData (GPG ASCII-Armored PEM) and the metadata API. There are a few optional fields that can be used modify the behavior of the package.
func New(keyData []byte, dataAPI metadata.API, opts ...Option) (ProviderKeyVerifier, error) {
	gpgVerifier, err := gpg_key_verifier.New(keyData)
	if err != nil {
		return nil, fmt.Errorf("cannot construct GPG key verifier: %w", err)
	}

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	// Default fields
	providerKeyVerifier := &providerKeyVerifier{
		gpgVerifier:     gpgVerifier,
		dataAPI:         dataAPI,
		httpClient:      httpClient,
		logger:          slog.New(slog.NewTextHandler(os.Stdout, nil)),
		versionsToCheck: 10,
		checkFn:         process,
	}

	for _, opt := range opts {
		opt(providerKeyVerifier)
	}

	return providerKeyVerifier, nil
}

// Option is used for providing options to New without changing the signature of New.
type Option func(c *providerKeyVerifier)
type CheckFn func(pkv providerKeyVerifier, ctx context.Context, version provider.Version) error

// WithVersionsToCheck is a functional option to set the number of versions to check for a provider.
func WithVersionsToCheck(versionsToCheck uint8) Option {
	return func(c *providerKeyVerifier) {
		c.versionsToCheck = versionsToCheck
	}
}

// WithLogger is a functional option to set the logger
func WithLogger(logger *slog.Logger) Option {
	return func(c *providerKeyVerifier) {
		c.logger = logger
	}
}

// WithHTTPClient is a functional option to set the http Client
func WithHTTPClient(httpClient http.Client) Option {
	return func(c *providerKeyVerifier) {
		c.httpClient = httpClient
	}
}

// WithCheckFn is a functional option to set the function used to check the provider version
func WithCheckFn(checkFn CheckFn) Option {
	return func(c *providerKeyVerifier) {
		c.checkFn = checkFn
	}
}

type providerKeyVerifier struct {
	gpgVerifier     gpg_key_verifier.GPGKeyVerifier
	dataAPI         metadata.API
	httpClient      http.Client
	versionsToCheck uint8
	logger          *slog.Logger
	checkFn         CheckFn
}
