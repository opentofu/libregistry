// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package common

// PackageTarget is a string that describes a platform and architecture combination, separated
// by an underscore (_).
//
// example: darwin_amd64
// swagger:model RegistryPackageTarget
type PackageTarget string
