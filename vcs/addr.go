package vcs

import (
	"regexp"
)

// nameRe is a common understanding of acceptable repository addresses. This may need to be changed later if it turns
// out that other VCS' support different address styles.
var nameRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// OrganizationAddr refers to an organization within a VCS system.
type OrganizationAddr struct {
	// Org is the URL fragment of an organization.
	Org string
}

func (o OrganizationAddr) String() string {
	return o.Org
}

// RepositoryAddr holds a reference to a repository. For simplicity, the current system does not support more complex
// URL structures.
type RepositoryAddr struct {
	OrganizationAddr
	// Name is the URL fragment of a repository.
	Name string
}

func (r RepositoryAddr) String() string {
	return r.OrganizationAddr.String() + "/" + r.Name
}

func (r RepositoryAddr) Validate() error {
	if !nameRe.MatchString(r.Org) {
		return &InvalidRepositoryAddrError{
			RepositoryAddr: r.String(),
		}
	}
	if !nameRe.MatchString(r.Name) {
		return &InvalidRepositoryAddrError{
			RepositoryAddr: r.String(),
		}
	}
	return nil
}
