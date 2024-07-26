// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package github

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/opentofu/libregistry/logger"
)

// Opt is a function that modifies the config.
type Opt func(config *Config) error

// Config holds the configuration for GitHub.
type Config struct {
	// Username to use for cloning in conjunction with a token.
	Username string
	// Token is the GitHub token to use when accessing the GitHub API and cloning.
	Token string
	// CheckoutRootDirectory is the root directory where repositories should be checked out. Defaults to the OS' temp
	// directory.
	CheckoutRootDirectory string
	// SkipCleanupWorkingCopyOnClose indicates that the working copy should not be cleaned up when it is closed.
	// Defaults to false, cleaning up the working copy.
	SkipCleanupWorkingCopyOnClose bool
	// GitPath holds the path to the git binary. Defaults to looking up the "git" or "git.exe" binaries in the path.
	GitPath string

	// Logger holds the logger to write any logs to.
	Logger logger.Logger
	// HTTPClient holds the HTTP client to use for API requests. Note that this only affects API and RSS feed requests,
	// but not git clone commands as those are done using the command line.
	HTTPClient *http.Client
}

// ApplyDefaults adds the default values if none are present.
func (c *Config) ApplyDefaults() {
	if c.CheckoutRootDirectory == "" {
		c.CheckoutRootDirectory = os.TempDir()
	}

	if c.GitPath == "" {
		c.GitPath = defaultGitPath
	}

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

// WithUsername sets the GitHub username to use for cloning a repository in conjunction with a token.
func WithUsername(username string) Opt {
	return func(config *Config) error {
		config.Username = username
		return nil
	}
}

// WithToken sets the GitHub token to use for authentication against the API when needed. If the username is also set,
// this token will be used for cloning a repository when needed.
func WithToken(token string) Opt {
	return func(config *Config) error {
		config.Token = token
		return nil
	}
}

// WithCheckoutRootDirectory sets a directory to use for repository checkouts.
func WithCheckoutRootDirectory(rootDir string) Opt {
	return func(config *Config) error {
		stat, err := os.Stat(rootDir)
		if err != nil {
			return fmt.Errorf("unusable checkout root directory (%w)", err)
		}
		if !stat.IsDir() {
			return fmt.Errorf("unusable checkout root directory (not a directory)")
		}
		rootDir, err = filepath.Abs(rootDir)
		if err != nil {
			return fmt.Errorf("failed to determine absolute path for %s (%v)", rootDir, err)
		}
		config.CheckoutRootDirectory = rootDir
		return nil
	}
}

// WithSkipCleanupWorkingCopyOnClose skips cleaning up the working directory when it is closed. This is useful when
// wanting to re-use the working directory and skip re-cloning the repository.
func WithSkipCleanupWorkingCopyOnClose(skip bool) Opt {
	return func(config *Config) error {
		config.SkipCleanupWorkingCopyOnClose = skip
		return nil
	}
}

// WithGitPath sets the path to the Git binary. Defaults to looking up the "git" or "git.exe" binaries in the path.
func WithGitPath(path string) Opt {
	return func(config *Config) error {
		cmd := exec.Command(path, "version")
		if err := cmd.Run(); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if exitErr.ExitCode() != 0 {
					return fmt.Errorf("git binary %s is not usable (git version exited with %d)", path, exitErr.ExitCode())
				}
			} else {
				return fmt.Errorf("git binary %s is not usable (%w)", path, err)
			}
		}
		config.GitPath = path
		return nil
	}
}

// WithLogger sets a logger to use for writing trace and debug information.
func WithLogger(logger logger.Logger) Opt {
	return func(config *Config) error {
		config.Logger = logger.WithName("GitHub")
		return nil
	}
}

// WithHTTPClient sets an HTTP client to use for API and RSS queries.
func WithHTTPClient(client *http.Client) Opt {
	return func(config *Config) error {
		config.HTTPClient = client
		return nil
	}
}
