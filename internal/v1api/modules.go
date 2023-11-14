package v1api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"registry-stable/internal"
	"registry-stable/internal/module"
)

// GenerateModuleResponses generates the response for the module version listing API endpoints.
// For more information see
// https://opentofu.org/docs/internals/module-registry-protocol/#list-available-versions-for-a-specific-module
// https://opentofu.org/docs/internals/module-registry-protocol/#download-source-code-for-a-specific-module-version
func (g Generator) GenerateModuleResponses(_ context.Context, namespace string, name string, targetSystem string) error {
	logger := slog.With(slog.String("namespace", namespace), slog.String("name", name), slog.String("targetSystem", targetSystem))

	// TODO: Get path calculation from somewhere else
	path := filepath.Join(namespace[0:1], namespace, name, targetSystem+".json")

	metadata, err := g.readModuleMetadata(path, logger)
	if err != nil {
		return err
	}

	// Right now the format is pretty much identical, however if we want to extend the results in the future to include
	// more information, we can do that here. (i.e. the root module, or the submodules)
	versionsResponse := make([]VersionResponseItem, len(metadata.Versions))
	for i, m := range metadata.Versions {
		versionsResponse[i] = VersionResponseItem{Version: m.Version}

		err := g.writeModuleVersionDownload(namespace, name, targetSystem, m.Version)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file for version %s: %w", m.Version, err)
		}
		logger.Debug("Wrote metadata version download file", slog.String("version", m.Version))
	}

	// Write the /versions response
	err = g.writeModuleVersionListing(namespace, name, targetSystem, versionsResponse)
	if err != nil {
		return err
	}

	return nil
}

// readModuleMetadata reads the module metadata file from the filesystem directly. This data should be the data fetched from the git repository.
func (g Generator) readModuleMetadata(path string, logger *slog.Logger) (*module.MetadataFile, error) {
	// list directories at the root of the fs
	dirs, err := fs.ReadDir(g.ModuleFS, ".")
	if err != nil {
		slog.Error("Failed to list directories", slog.Any("err", err))
		os.Exit(1)
	}

	for _, d := range dirs {
		slog.Info("Found directory", slog.String("dir", d.Name()))
	}

	// open the file
	metadataFile, err := fs.ReadFile(g.ModuleFS, path)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}

	// Read the file contents into a Module[] struct
	var metadata module.MetadataFile
	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		return nil, err
	}

	logger.Debug("Loaded Modules", slog.Any("count", len(metadata.Versions)))

	return &metadata, nil
}

// writeModuleVersionListing writes the file containing the module version listing.
// This data  is to be consumed when an end user requests /v1/modules/{namespace}/{name}/{targetSystem}/versions
func (g Generator) writeModuleVersionListing(namespace string, name string, targetSystem string, versions []VersionResponseItem) error {
	destinationDir := filepath.Join(g.DestinationDir, "v1", "modules", namespace, name, targetSystem)
	if err := g.FileWriter.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(destinationDir, "versions")

	marshalled, err := json.Marshal(ModuleVersionListingResponse{Modules: []ModuleVersionListingResponseItem{{Versions: versions}}})
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	err = g.FileWriter.WriteFile(filePath, marshalled, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// writeModuleVersionDownload writes the file containing the download link for the module version.
// This data is to be consumed when an end user requests /v1/modules/{namespace}/{name}/{targetSystem}/{version}/download
func (g Generator) writeModuleVersionDownload(namespace string, name string, system string, version string) interface{} {
	// the file should just contain a link to GitHub to download the tarball, ie:
	// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=v5.30.0
	contents := fmt.Sprintf("git::github.com/%s/terraform-%s-%s?ref=%s", namespace, name, system, version)

	// trim the v from the version
	ver := internal.TrimTagPrefix(version)

	destinationDir := filepath.Join(g.DestinationDir, "v1", "modules", namespace, name, system, ver)
	if err := g.FileWriter.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(destinationDir, "download")
	err := g.FileWriter.WriteFile(filePath, []byte(contents), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
