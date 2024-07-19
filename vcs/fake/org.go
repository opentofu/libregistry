package fake

import (
	"github.com/opentofu/libregistry/vcs"
)

type org struct {
	users        map[string]struct{}
	repositories map[vcs.RepositoryAddr]*repository
}
