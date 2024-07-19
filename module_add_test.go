package libregistry_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/opentofu/libregistry"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/memory"
	"github.com/opentofu/libregistry/vcs/github"
)

func ExampleAPI_AddModule() {
	ghClient, err := github.New(os.Getenv("GITHUB_TOKEN"), nil)
	if err != nil {
		panic(err)
	}

	storage := memory.New()

	dataAPI, err := metadata.New(storage)
	if err != nil {
		panic(err)
	}

	registry, err := libregistry.New(
		ghClient,
		dataAPI,
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	if err := registry.AddModule(ctx, "terraform-aws-modules/terraform-aws-iam"); err != nil {
		panic(err)
	}

	// Manually read the registry data file:
	jsonFile, err := storage.GetFile(ctx, "modules/t/terraform-aws-modules/iam/aws.json")
	if err != nil {
		panic(err)
	}

	var data map[string]any
	if err := json.Unmarshal(jsonFile, &data); err != nil {
		panic(err)
	}
	fmt.Printf("Latest version: %s", data["versions"].([]any)[0].(map[string]any)["version"].(string))
}
