// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package fakevcs

import (
	"io/fs"
	"time"

	"github.com/opentofu/libregistry/vcs"
)

// NewWithOpts creates a fake, in-memory VCSClient implementation for testing use with added options.
func NewWithOpts(opts ...Opt) (VCSClient, error) {
	cfg := Config{}

	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	cfg.ApplyDefaults()

	return &inMemoryVCS{
		config:        cfg,
		users:         map[vcs.Username]struct{}{},
		organizations: map[vcs.OrganizationAddr]*org{},
	}, nil
}

// New creates a fake, in-memory VCSClient implementation for testing use.
func New() VCSClient {
	client, err := NewWithOpts()
	if err != nil {
		panic(err)
	}
	return client
}

type Opt func(config *Config) error

func WithTimeSource(timeSource func() time.Time) Opt {
	return func(config *Config) error {
		config.TimeSource = timeSource
		return nil
	}
}

type Config struct {
	TimeSource func() time.Time
}

func (c *Config) ApplyDefaults() {
	if c.TimeSource == nil {
		c.TimeSource = time.Now
	}
}

type VCSClient interface {
	vcs.Client

	CreateOrganization(organization vcs.OrganizationAddr) error
	CreateRepository(repository vcs.RepositoryAddr, info vcs.RepositoryInfo) error
	CreateVersion(repository vcs.RepositoryAddr, version vcs.VersionNumber, content fs.ReadDirFS) error
	AddAsset(repository vcs.RepositoryAddr, version vcs.VersionNumber, name vcs.AssetName, data []byte) error
	AddUser(username vcs.Username) error
	AddMember(organization vcs.OrganizationAddr, username vcs.Username) error
}
