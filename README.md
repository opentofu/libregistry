# Go library for the OpenTofu registry

This Go library implements the OpenTofu registry and also provides a library to access the underlying data. You can install this library by running:

```
go get github.com/opentofu/libregistry
```

## The metadata API

The metadata API is a low-level API that allows you to access the stored registry data structure in the [registry repository](https://github.com/opentofu/registry). Use this API when you need to work with the registry data without needing online functions, such as refreshing a module or provider.

You can use the metadata API like this:

```go
package main

import (
    "context"

    "github.com/opentofu/libregistry/metadata"
    "github.com/opentofu/libregistry/metadata/storage/filesystem"
)

func main() {
    metadataAPI, err := metadata.New(filesystem.New("path/to/registry/data"))
    if err != nil {
        panic(err)
    }
    modules, err := metadataAPI.ListModules(context.Background())
    if err != nil {
        panic(err)
    }
    
    // Do something with modules here.
}
```

## The registry API

The `libregistry` package contains the top level registry API. It implements the functions that are triggered from GitHub Actions, such as adding a module, etc.

You can use the registry API like this:

```go
package main

import (
	"context"
	"os"

	"github.com/opentofu/libregistry"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/filesystem"
	"github.com/opentofu/libregistry/vcs/github"
)

func main() {
	ghClient, err := github.New(os.Getenv("GITHUB_TOKEN"), nil, nil)
	if err != nil {
		panic(err)
	}

	storage := filesystem.New("path/to/registry/data")

	metadataAPI, err := metadata.New(storage)
	if err != nil {
		panic(err)
	}

	registry, err := libregistry.New(
		ghClient,
		metadataAPI,
	)
	if err != nil {
		panic(err)
	}

	if err := registry.AddModule(context.TODO(), "terraform-aws-modules/terraform-aws-iam"); err != nil {
		panic(err)
	}
}
```

## VCS implementations

This library supports pluggable VCS systems. We run on GitHub by default, but you may be interested in implementing a VCS backend for a different system. Check out the [vcs](vcs) package for the VCS interface. Note, that the implementation still assumes that you will have an organization/repository structure and many systems, such as the registry UI, still assume that the VCS system will be git.

## Metadata storage

You may also be interested in storing the metadata somewhere else than the local filesystem. For this purpose, check out the [metadata/storage](metadata/storage) package, which contains the interface for defining storages.