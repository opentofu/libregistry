// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package logger_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/opentofu/libregistry/logger"
)

func TestWriter(t *testing.T) {
	ctx := context.Background()
	log := &collectingTestLogger{}

	const testMessage1 = "Hello world!"
	const testMessage2 = "Hello world2"
	const testMessage3 = "Hello world3"

	wr := logger.NewWriter(ctx, log, logger.LevelDebug, "")

	if _, err := wr.Write([]byte(testMessage1)); err != nil {
		t.Fatal(err)
	}
	if len(log.lines) != 0 {
		t.Fatalf("Incorrect number of lines logged: %d", len(log.lines))
	}

	if _, err := wr.Write([]byte("\n")); err != nil {
		t.Fatal(err)
	}
	if len(log.lines) != 1 {
		t.Fatalf("Incorrect number of lines logged: %d", len(log.lines))
	}
	if log.lines[0] != testMessage1 {
		t.Fatalf("Incorrect log message: %s", log.lines[0])
	}

	if _, err := wr.Write([]byte(testMessage2 + "\n" + testMessage3)); err != nil {
		t.Fatal(err)
	}
	if len(log.lines) != 2 {
		t.Fatalf("Incorrect number of lines logged: %d", len(log.lines))
	}
	if log.lines[1] != testMessage2 {
		t.Fatalf("Incorrect log message: %s", log.lines[1])
	}

	if err := wr.Close(); err != nil {
		t.Fatal(err)
	}
	if len(log.lines) != 3 {
		t.Fatalf("Incorrect number of lines logged: %d", len(log.lines))
	}
	if log.lines[2] != testMessage3 {
		t.Fatalf("Incorrect log message: %s", log.lines[2])
	}
}

type collectingTestLogger struct {
	lines []string
}

func (t *collectingTestLogger) WithName(name string) logger.Logger {
	return t
}

func (t *collectingTestLogger) log(message string, args []any) {
	t.lines = append(t.lines, fmt.Sprintf(message, args...))
}

func (t *collectingTestLogger) Trace(ctx context.Context, message string, args ...any) {
	t.log(message, args)
}

func (t *collectingTestLogger) Debug(ctx context.Context, message string, args ...any) {
	t.log(message, args)
}

func (t *collectingTestLogger) Info(ctx context.Context, message string, args ...any) {
	t.log(message, args)
}

func (t *collectingTestLogger) Warn(ctx context.Context, message string, args ...any) {
	t.log(message, args)
}

func (t *collectingTestLogger) Error(ctx context.Context, message string, args ...any) {
	t.log(message, args)
}
