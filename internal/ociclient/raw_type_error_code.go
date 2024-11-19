// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

type OCIRawErrorCode string

const (
	OCIErrorCodeBlobUnknown         OCIRawErrorCode = "BLOB_UNKNOWN"
	OCIErrorCodeUploadInvalid       OCIRawErrorCode = "BLOB_UPLOAD_INVALID"
	OCIErrorCodeUploadUnknown       OCIRawErrorCode = "BLOB_UPLOAD_UNKNOWN"
	OCIErrorCodeDigestInvalid       OCIRawErrorCode = "DIGEST_INVALID"
	OCIErrorCodeManifestBlobUnknown OCIRawErrorCode = "MANIFEST_BLOB_UNKNOWN"
	OCIErrorCodeManifestInvalid     OCIRawErrorCode = "MANIFEST_INVALID"
	OCIErrorCodeManifestUnknown     OCIRawErrorCode = "MANIFEST_UNKNOWN"
	OCIErrorCodeNameInvalid         OCIRawErrorCode = "NAME_INVALID"
	OCIErrorCodeNameUnknown         OCIRawErrorCode = "NAME_UNKNOWN"
	OCIErrorCodeSizeInvalid         OCIRawErrorCode = "SIZE_INVALID"
	OCIErrorCodeUnauthorized        OCIRawErrorCode = "UNAUTHORIZED"
	OCIErrorCodeDenied              OCIRawErrorCode = "DENIED"
	OCIErrorCodeUnsupported         OCIRawErrorCode = "UNSUPPORTED"
	OCIErrorCodeTooManyRequests     OCIRawErrorCode = "TOOMANYREQESTS"
)
