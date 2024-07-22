// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package module

// Metadata represents all the metadata for a module. This includes the list of versions available for the module.
// This structure represents the file in modules/o/opentofu/somemodule/platform.json.
type Metadata struct {
	// Versions lists all available versions of a Namespace-Name-TargetSystem combination.
	Versions VersionList `json:"versions"`
}
