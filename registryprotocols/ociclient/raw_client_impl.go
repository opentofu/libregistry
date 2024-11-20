// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentofu/libregistry/logger"
	"net/http"
	"strings"
)

type rawClient struct {
	httpClient  *http.Client
	credentials ScopedCredentials
	logger      logger.Logger
}

func (r *rawClient) Check(
	ctx context.Context,
	registry OCIRegistry,
) (OCIWarnings, error) {
	if err := registry.Validate(); err != nil {
		return nil, err
	}

	r.logger.Debug(ctx, "Performing access check for registry %s...", registry)

	endpoint := fmt.Sprintf("https://%s/v2/", registry)

	response, warnings, err := getWithAuthentication(ctx, r, OCIAddr{
		Registry: registry,
		Name:     "",
	}, endpoint, nil)
	if err != nil {
		// Some OCI registries, such as ghcr.io, don't allow anonymous authentication for the /v2/ endpoint.
		var authRequiredError *OCIRawAuthenticationRequiredError
		if !errors.As(err, &authRequiredError) {
			var ociError *OCIRawErrors
			if errors.As(err, &ociError) {
				r.logger.Warn(ctx, "Access check for registry %s failed and no structured error returned, this may not be an OCI registry. (%v)", registry, err)
			} else {
				r.logger.Debug(ctx, "Access check for registry %s failed: %v", registry, err)
			}
			return warnings, err
		}
		r.logger.Debug(ctx, "Access check for registry %s failed with an authentication requirement, ignoring error: %v", registry, err)
	}
	if response != nil {
		_ = response.Body.Close()
	}
	return warnings, nil
}

func (r *rawClient) ContentDiscovery(
	ctx context.Context,
	addr OCIAddr,
) (OCIRawContentDiscoveryResponse, OCIWarnings, error) {
	if err := addr.Validate(); err != nil {
		return OCIRawContentDiscoveryResponse{}, nil, err
	}

	r.logger.Trace(ctx, "Performing content discovery for %s...", addr)

	var result OCIRawContentDiscoveryResponse
	endpoint := fmt.Sprintf("https://%s/v2/%s/tags/list", addr.Registry, addr.Name)
	response, warnings, err := getWithAuthentication(ctx, r, addr, endpoint, []string{"application/json"})
	if err != nil {
		return result, warnings, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&result); err != nil {
		return result, warnings, newInvalidOCIResponseError(
			"Failed to decode OCI content discovery response",
			err,
		)
	}
	return result, warnings, err

}

func (r *rawClient) GetManifest(
	ctx context.Context,
	addrRef OCIAddrWithReference,
) (OCIRawManifest, OCIWarnings, error) {
	if err := addrRef.Validate(); err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("https://%s/v2/%s/manifests/%s", addrRef.Registry, addrRef.Name, addrRef.Reference)
	accept := []string{
		string(MediaTypeDockerImage),
		string(MediaTypeDockerImageList),
		string(MediaTypeOCIImage),
		string(MediaTypeOCIImageIndex),
	}
	response, warnings, err := getWithAuthentication(ctx, r, addrRef.OCIAddr, endpoint, accept)
	if err != nil {
		return nil, warnings, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var result OCIRawManifest
	contentType := response.Header.Get("Content-Type")
	switch OCIRawMediaType(contentType) {
	case MediaTypeOCIImageIndex:
		fallthrough
	case MediaTypeDockerImageList:
		result = &OCIRawImageIndexManifest{}
	case MediaTypeOCIImage:
		fallthrough
	case MediaTypeDockerImage:
		result = &OCIRawImageManifest{}
	default:
		return nil, warnings, fmt.Errorf("protocol error: the OCI registry server ignored the Accept header")
	}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&result); err != nil {
		return result, warnings, newInvalidOCIResponseError(
			"Failed to decode OCI manifest response",
			err,
		)
	}
	return result, warnings, err
}

func (r *rawClient) GetBlob(
	ctx context.Context,
	addrDigest OCIAddrWithDigest,
) (OCIRawBlob, OCIWarnings, error) {
	if err := addrDigest.Validate(); err != nil {
		return OCIRawBlob{}, nil, err
	}
	endpoint := fmt.Sprintf("https://%s/v2/%s/blobs/%s", addrDigest.Registry, addrDigest.Name, addrDigest.Digest)
	accept := []string{
		"application/octet-stream",
	}
	response, warnings, err := getWithAuthentication(ctx, r, addrDigest.OCIAddr, endpoint, accept)
	if err != nil {
		return OCIRawBlob{}, warnings, err
	}
	return OCIRawBlob{
		response.Body,
		OCIRawMediaType(response.Header.Get("Content-Type")),
	}, warnings, nil
}

func getWithAuthentication(
	ctx context.Context,
	r *rawClient,
	addr OCIAddr,
	endpoint string,
	accept []string,
) (*http.Response, OCIWarnings, error) {
	response, warnings, err := tryRequest(ctx, r, addr, endpoint, accept)
	if err == nil {
		return response, warnings, err
	}

	var authRequired *OCIRawAuthenticationRequiredError
	if !errors.As(err, &authRequired) {
		// Not an authentication required error.
		r.logger.Trace(ctx, "The request to %s failed with the following error: %v", endpoint, err)
		return response, warnings, err
	}
	// Authentication required, authenticate and try again:

	// We now have to try and authenticate against the realm endpoint
	authSchemes := authRequired.GetAuthSchemes("Bearer")
	if len(authSchemes) == 0 {
		r.logger.Debug(ctx, "The request to %s failed with an authentication failure, but no suitable (Bearer) authentication endpoints were provided.", endpoint, err)
		return response, warnings, err
	}
	r.logger.Debug(ctx, "The request to %s failed with an authentication failure, attempting to obtain a bearer token...", endpoint)
	for _, authScheme := range authSchemes {
		realm, ok := authScheme.GetParam("realm")
		if !ok || realm == "" {
			continue
		}
		queryString := authScheme.ParamsAsQueryString("realm")
		if strings.Contains(realm, "?") {
			realm += "&" + queryString
		} else {
			realm += "?" + queryString
		}

		// Find credentials:
		var creds *ClientCredentials
		filter := func(scope OCIScope, creds *ClientCredentials) bool {
			return creds != nil && creds.Basic != nil
		}
		if addr.Name == "" {
			// This is for the Check() call only where no name is set.
			creds = r.credentials.GetCredentialsForRegistry(addr.Registry, filter)
		} else {
			creds = r.credentials.GetCredentialsForAddr(addr, filter)
		}
		authorization := ""
		if creds != nil {
			authorization = "Basic " + base64.RawURLEncoding.EncodeToString([]byte(creds.Basic.Username+":"+creds.Basic.Password))
		}
		response, newWarnings, err := getRequest(ctx, r, realm, []string{"application/json"}, authorization)
		if err != nil {
			r.logger.Debug(ctx, "Authentication request to %s for endpoint %s failed, trying next authentication method (%v).", realm, endpoint, err)
			continue
		}

		authResponse := ClientBearerTokenCredentials{}
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&authResponse); err != nil {
			_ = response.Body.Close()
			r.logger.Debug(ctx, "Protocol error: authentication request to %s for endpoint %s returned an invalid response, trying next authentication method... (%v)", realm, endpoint, err)
			continue
		}
		_ = response.Body.Close()

		if creds == nil {
			creds = &ClientCredentials{}
			r.credentials.SetCredentials(addr.ToScope(), creds)
		}
		// We are storing the credentials we just obtained so we don't have to re-authenticate again.
		creds.Bearer = &authResponse
		warnings = append(warnings, newWarnings...)
		r.logger.Debug(ctx, "Authentication for %s successful on realm %s, retrying previous request.", endpoint, realm)

		response, newWarnings, err = tryRequest(ctx, r, addr, endpoint, accept)
		warnings = append(warnings, newWarnings...)
		if err == nil {
			// We don't close the response on purpose so the caller can use it.
			return response, warnings, err
		}
		r.logger.Debug(ctx, "Authentication for %s failed on realm %s, trying next authentication method... (%v)", endpoint, realm, err)
		if response != nil {
			_ = response.Body.Close()
		}
	}
	r.logger.Trace(ctx, "No more bearer authentication methods remaining, authentication for %s failed.", endpoint)
	return response, warnings, authRequired
}

func tryRequest(ctx context.Context, r *rawClient, addr OCIAddr, endpoint string, accept []string) (*http.Response, OCIWarnings, error) {
	filter := func(scope OCIScope, creds *ClientCredentials) bool {
		return creds != nil && creds.Bearer != nil
	}
	var creds *ClientCredentials
	if addr.Name == "" {
		// We are handling a Check() call where there is no name
		creds = r.credentials.GetCredentialsForRegistry(addr.Registry, filter)
	} else {
		creds = r.credentials.GetCredentialsForAddr(addr, filter)
	}
	authorization := ""
	if creds != nil {
		authorization = "Bearer " + creds.Bearer.Token
	}
	return getRequest(ctx, r, endpoint, accept, authorization)
}

func getRequest(
	ctx context.Context,
	r *rawClient,
	endpoint string,
	accept []string,
	authorization string,
) (*http.Response, OCIWarnings, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		// TODO typed error
		return nil, nil, fmt.Errorf("failed to construct HTTP request (%w)", err)
	}

	if len(accept) > 0 {
		req.Header.Add("Accept", strings.Join(accept, ", "))
	}
	if len(authorization) > 0 {
		r.logger.Trace(ctx, "Sending request to %s with credentials...", endpoint)
		req.Header.Add("Authorization", authorization)
	} else {
		r.logger.Trace(ctx, "Sending request to %s without credentials...", endpoint)
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request to %s (%w)", endpoint, err)
	}

	warnings := OCIWarnings(resp.Header.Values("Warning"))
	for _, warning := range warnings {
		r.logger.Trace(ctx, "OCI endpoint %s returned a warning: %s", endpoint, warning)
	}

	e := &OCIRawErrors{}
	switch {
	case resp.StatusCode > 199 && resp.StatusCode < 300:
		// We don't close the body here so the caller can use it.
		r.logger.Trace(ctx, "Request to %s successful.", endpoint)
		return resp, warnings, nil
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		// We treat the 401 and 403 the same, even though it's not entirely RFC-conformant because some registries,
		// such as ghcr.io, use both error codes in certain situations incorrectly.
		defer func() {
			_ = resp.Body.Close()
		}()
		var authSchemes []OCIRawAuthScheme
		for _, wwwAuthenticateHeader := range resp.Header.Values("WWW-Authenticate") {
			// Note: according to RFC 7235 multiple WWW-Authenticate headers may be present.
			schemes, err := parseWWWAuthenticate(wwwAuthenticateHeader)
			if err != nil {
				// Invalid www-authenticate header, the response is malformed and we can't use it
				r.logger.Warn(ctx, "The registry at %s returned a malformed WWW-Authenticate header. (%v)", endpoint, err)
				return nil, warnings, fmt.Errorf("cannot decode www-authenticate header from OCI registry (%w)", err)
			}
			authSchemes = append(authSchemes, schemes...)
		}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(e); err != nil {
			// We can't decode the error response properly, return the authentication required response
			// without a cause.
			r.logger.Trace(ctx, "Cannot decode JSON error response from %s. (%v)", endpoint, err)
		}
		r.logger.Trace(ctx, "Request to %s returned an authentication failure.", endpoint)
		return nil, warnings, newOCIRawAuthenticationRequiredError(
			endpoint,
			authSchemes,
			e,
		)
	default:
		defer func() {
			_ = resp.Body.Close()
		}()
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(e); err != nil {
			r.logger.Trace(ctx, "The endpoint at %s returned an unexpected HTTP error %d. Additionally, the error could not be interpreted as an OCI error response. (%v)", endpoint, resp.StatusCode, err)
			return nil, warnings, fmt.Errorf("cannot decode OCI response from %s into %T (%w)", endpoint, e, err)
		}
		r.logger.Trace(ctx, "The endpoint at %s returned an unexpected HTTP error %d. (%v)", endpoint, resp.StatusCode, e)
		return nil, warnings, e
	}
}

var _ RawOCIClient = &rawClient{}
