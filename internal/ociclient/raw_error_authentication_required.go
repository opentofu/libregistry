// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"fmt"
	"net/url"
	"strings"
)

type OCIRawAuthScheme struct {
	Type   string
	Params map[string]string
}

func (o OCIRawAuthScheme) GetParam(name string) (string, bool) {
	for k, v := range o.Params {
		if strings.ToLower(k) == strings.ToLower(name) {
			return v, true
		}
	}
	return "", false
}

func (o OCIRawAuthScheme) ParamsAsQueryString(withoutParam string) string {
	var parts []string
	for k, v := range o.Params {
		if withoutParam != "" && strings.ToLower(k) == strings.ToLower(withoutParam) {
			continue
		}
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	return strings.Join(parts, "&")
}

// OCIRawAuthenticationRequiredError indicates that authentication is required to access a specific endpoint.
type OCIRawAuthenticationRequiredError struct {
	Endpoint    string
	AuthSchemes []OCIRawAuthScheme
	Cause       error
}

// GetAuthSchemes returns the authentication schemes from the WWW-Authenticate headers matching the given authScheme.
func (o OCIRawAuthenticationRequiredError) GetAuthSchemes(authScheme string) []OCIRawAuthScheme {
	var result []OCIRawAuthScheme
	for _, scheme := range o.AuthSchemes {
		if strings.ToLower(scheme.Type) == strings.ToLower(authScheme) {
			result = append(result, scheme)
		}
	}
	return result
}

func (o OCIRawAuthenticationRequiredError) Error() string {
	if o.Cause != nil {
		return fmt.Sprintf("Authentication required while accessing %s (%v)", o.Endpoint, o.Cause)
	}
	return fmt.Sprintf("Authentication required while accessing %s", o.Endpoint)
}

func (o OCIRawAuthenticationRequiredError) Unwrap() error {
	return o.Cause
}

func newOCIRawAuthenticationRequiredError(endpoint string, authSchemes []OCIRawAuthScheme, cause error) error {
	return &OCIRawAuthenticationRequiredError{
		Endpoint:    endpoint,
		AuthSchemes: authSchemes,
		Cause:       cause,
	}
}
