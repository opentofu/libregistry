package github

import (
	"context"
	"strings"

	"github.com/opentofu/libregistry/vcs"
)

// New creates a new GitHub VCS client.
func New(token string) (vcs.Client, error) {
	return &github{
		token: token,
	}, nil
}

type github struct {
	token string
}

func (g github) ParseRepositoryAddr(ref string) (vcs.RepositoryAddr, error) {
	ref = strings.TrimPrefix(ref, "github.com/")
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) != 2 {
		return vcs.RepositoryAddr{}, &vcs.InvalidRepositoryAddrError{
			RepositoryAddr: ref,
		}
	}
	result := vcs.RepositoryAddr{
		OrganizationAddr: vcs.OrganizationAddr{Org: parts[0]},
		Name:             parts[1],
	}
	return result, result.Validate()
}

func (g github) ListVersions(ctx context.Context, repository vcs.RepositoryAddr) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (g github) ListAssets(ctx context.Context, repository vcs.RepositoryAddr, version string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (g github) DownloadAsset(ctx context.Context, repository vcs.RepositoryAddr, version string, asset string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g github) HasPermission(ctx context.Context, organization vcs.OrganizationAddr) (bool, error) {
	//TODO implement me
	panic("implement me")
}
