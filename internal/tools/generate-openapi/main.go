// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 2 {
		_, _ = os.Stderr.Write([]byte("Usage: generate-openapi github.com/opentofu/libregistry/..."))
		os.Exit(1)
	}
	pkg := os.Args[1]
	cmd := exec.Command("go", "run", "github.com/go-swagger/go-swagger/cmd/swagger@v0.31.0", "generate", "spec", "-o", "openapi.yml", "-m", "--include=^(github.com/opentofu/libregistry/registry/common|"+pkg+")$")
	cmd.Env = append(os.Environ(), "SWAGGER_GENERATE_EXTENSION=false")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
			return
		}
		log.Print(err)
		os.Exit(1)
	}
}
