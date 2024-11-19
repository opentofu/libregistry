// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

// ScopedCredentials describes a configuration structure that holds credentials on a per-registry basis.
// The key may contain a hostname, or it may contain a hostname and a name for which the credentials apply.
type ScopedCredentials map[OCIScopeString]*ClientCredentials

func (s ScopedCredentials) Validate() error {
	for scopeString, creds := range s {
		scope, err := scopeString.Parse()
		if err != nil {
			return newConfigurationError("invalid registry scope: "+string(scopeString), err)
		}
		if err := scope.Validate(); err != nil {
			return newConfigurationError("invalid registry scope: "+scope.String(), err)
		}
		if err := creds.Validate(); err != nil {
			return newConfigurationError("invalid credentials configured for "+scope.String(), err)
		}
	}
	return nil
}

func (s ScopedCredentials) SetCredentials(scope OCIScope, creds *ClientCredentials) {
	s[scope.ScopeString()] = creds
}

func (s ScopedCredentials) GetCredentialsForRegistry(registry OCIRegistry, filter func(scope OCIScope, creds *ClientCredentials) bool) *ClientCredentials {
	for scopeString, creds := range s {
		scope, err := scopeString.Parse()
		if err != nil {
			continue
		}
		if scope.MatchesRegistry(registry) && (filter == nil || filter(scope, creds)) {
			return creds
		}
	}
}

// GetCredentialsForAddr returns the most appropriate credentials for the given scope.
func (s ScopedCredentials) GetCredentialsForAddr(addr OCIAddr, filter func(scope OCIScope, creds *ClientCredentials) bool) *ClientCredentials {
	for scopeString, creds := range s {
		scope, err := scopeString.Parse()
		if err != nil {
			continue
		}
		if scope.MatchesAddr(addr, true) && (filter == nil || filter(scope, creds)) {
			return creds
		}
	}
	for scopeString, creds := range s {
		scope, err := scopeString.Parse()
		if err != nil {
			continue
		}
		if scope.MatchesAddr(addr, false) && (filter == nil || filter(scope, creds)) {
			return creds
		}
	}
	return nil
}

var _ validatable = ScopedCredentials{}
