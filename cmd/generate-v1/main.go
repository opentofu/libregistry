package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"registry-stable/internal/provider"
	"registry-stable/internal/repository-metadata-files/module"
	"registry-stable/internal/v1api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger) // TODO REMOVE ME
	logger.Info("Generating v1 API responses")

	ctx := context.Background()

	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	destinationDir := flag.String("destination", "../generated", "Directory to write the generated responses to")

	flag.Parse()

	v1APIGenerator := v1api.Generator{
		DestinationDir: *destinationDir,

		ModuleDirectory:   *moduleDataDir,
		ProviderDirectory: *providerDataDir,
	}

	err := v1APIGenerator.WriteWellKnownFile(ctx)
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	modules, err := module.ListModules(*moduleDataDir)
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	for _, m := range modules {
		logger.Info("Generating", slog.String("module", m.Namespace+"/"+m.Name+"/"+m.TargetSystem))
		err := v1APIGenerator.GenerateModuleResponses(ctx, m)
		if err != nil {
			logger.Error("Failed to generate module version listing response", slog.Any("err", err))
			os.Exit(1)
		}
		logger.Info("Generated", slog.String("module", m.Namespace+"/"+m.Name+"/"+m.TargetSystem))
	}

	providers, err := provider.ListProviders(*providerDataDir, slog.Default())
	if err != nil {
		slog.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}

	for _, p := range providers {
		err := v1APIGenerator.GenerateProviderResponses(ctx, p)
		if err != nil {
			p.Logger.Error("Failed to generate provider version listing response", slog.Any("err", err))
			os.Exit(1)
		}
	}

	slog.Info("Completed generating v1 API responses")
}
