// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package logger

import (
	"context"
	"testing"
)

// NewTestLogger produces a new logger that writes to a Go *testing.T.
func NewTestLogger(t *testing.T) Logger {
	return &testLogger{
		t, "",
	}
}

type testLogger struct {
	t      *testing.T
	prefix string
}

func (t testLogger) WithName(name string) Logger {
	return &testLogger{
		t.t,
		name + "\t",
	}
}

func (t testLogger) log(ctx context.Context, level string, message string, args ...any) {
	if t.prefix != "" {
		t.t.Logf(t.prefix+level+"\t"+message, args...)
	} else {
		t.t.Logf(level+"\t"+message, args...)
	}
}

func (t testLogger) Trace(ctx context.Context, message string, args ...any) {
	t.log(ctx, "TRACE", message, args...)
}

func (t testLogger) Debug(ctx context.Context, message string, args ...any) {
	t.log(ctx, "DEBUG", message, args...)
}

func (t testLogger) Info(ctx context.Context, message string, args ...any) {
	t.log(ctx, "INFO", message, args...)
}

func (t testLogger) Warn(ctx context.Context, message string, args ...any) {
	t.log(ctx, "WARN", message, args...)
}

func (t testLogger) Error(ctx context.Context, message string, args ...any) {
	t.log(ctx, "ERROR", message, args...)
}
