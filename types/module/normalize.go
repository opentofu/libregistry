package module

import (
	"strings"
)

func NormalizeNamespace(namespace string) string {
	return strings.ToLower(namespace)
}

func NormalizeName(name string) string {
	return strings.ToLower(name)
}

func NormalizeModuleTargetSystem(name string) string {
	return strings.ToLower(name)
}

func NormalizeModuleAddr(moduleAddr Addr) Addr {
	return moduleAddr.Normalize()
}