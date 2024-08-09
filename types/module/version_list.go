// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"slices"
)

// VersionList is a slice of versions.
//
// swagger:model ModuleVersionList
type VersionList []Version

// Merge merges the current list with another list and returns the new merged list.
func (v VersionList) Merge(other VersionList) VersionList {
	verSet := map[VersionNumber]Version{}
	for _, ver := range v {
		verSet[ver.Version.Normalize()] = ver
	}
	for _, ver := range other {
		verSet[ver.Version.Normalize()] = ver
	}
	newVersions := make(VersionList, len(verSet))
	i := 0
	for _, ver := range verSet {
		newVersions[i] = ver.Normalize()
		i++
	}
	newVersions.Sort()
	return newVersions
}

// Sort returns a sorted copy of the version list.
func (v VersionList) Sort() {
	semverSortFunc := func(a, b Version) int {
		return -a.Compare(b)
	}
	slices.SortFunc(v, semverSortFunc)
}

func (v VersionList) Equals(other VersionList) bool {
	if len(v) != len(other) {
		return false
	}
	for i := 0; i < len(v); i++ {
		if !v[i].Equals(other[i]) {
			return false
		}
	}
	return true
}
