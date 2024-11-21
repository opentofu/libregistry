// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ociclient

type validatable interface {
	Validate() error
}
