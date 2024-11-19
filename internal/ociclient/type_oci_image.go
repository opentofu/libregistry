// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "io"

// PulledOCIImage is an iterator interface to read pulled OCI images. You will have to call Next()
// and then read the file contents until it returns false. Make sure you only call Close()
// when you read all files from the image as this will clean up the temporary files.
//
// Developer note: this interface is intentionally created as an iterator since it allows us to
// read directly from the consecutive tar files returned from the OCI registry. It will also allow
// us to make layer downloads async or on-demand in the future if needed.
type PulledOCIImage interface {
	// Next moves the pointer to the next file in the image.
	Next() bool
	// Filename returns the filename of the current file in the image. Note that you are responsible
	// for validating if this filename is causing issues with an underlying filesystem or presents a
	// security risk.
	Filename() string

	io.ReadCloser
}
