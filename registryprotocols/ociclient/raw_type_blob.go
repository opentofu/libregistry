// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "io"

// OCIRawBlob contains a reader for a blob, as well as its content type.
type OCIRawBlob struct {
	io.ReadCloser

	ContentType OCIRawMediaType
}
