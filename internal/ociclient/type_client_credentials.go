// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// ClientCredentials holds the configuration of credentials for a registry. This may be either via basic
// auth or bearer tokens. The client MAY modify this if it obtains a bearer token based on the basic
// auth credentials. When both basic and bearer credentials are present, the bearer token is preferred.
type ClientCredentials struct {
	Basic  *ClientBasicCredentials       `json:"basic,omitempty"`
	Bearer *ClientBearerTokenCredentials `json:"bearer,omitempty"`
}

func (c ClientCredentials) Validate() error {
	if c.Basic == nil && c.Bearer == nil {
		return newConfigurationError("either basic or bearer authentication must be set", nil)
	}
	if c.Basic != nil {
		err := c.Basic.Validate()
		if err != nil {
			return newConfigurationError("invalid basic credentials", err)
		}
	}
	if c.Bearer != nil {
		err := c.Bearer.Validate()
		if err != nil {
			return newConfigurationError("invalid bearer token credentials", err)
		}
	}
	return nil
}

var _ validatable = ClientCredentials{}

// ClientBasicCredentials contains a username and password passed via the Authorization header.
type ClientBasicCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c ClientBasicCredentials) Validate() error {
	if c.Username == "" {
		return newConfigurationError(
			"the username cannot be empty", nil,
		)
	}
	return nil
}

var _ validatable = ClientBasicCredentials{}

// ClientBearerTokenCredentials contains a bearer token for direct authentication.
type ClientBearerTokenCredentials struct {
	Token string `json:"token"`
}

func (c ClientBearerTokenCredentials) Validate() error {
	if c.Token == "" {
		return newConfigurationError(
			"the token cannot be empty", nil,
		)
	}
	return nil
}

var _ validatable = ClientBearerTokenCredentials{}
