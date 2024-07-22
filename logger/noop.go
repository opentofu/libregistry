// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"context"
)

// NewNoopLogger creates a logger that does nothing.
func NewNoopLogger() Logger {
	return &noopLogger{}
}

type noopLogger struct {
}

func (n noopLogger) WithName(_ string) Logger {
	return n
}

func (n noopLogger) Trace(ctx context.Context, message string, args ...any) {

}

func (n noopLogger) Debug(ctx context.Context, message string, args ...any) {

}

func (n noopLogger) Info(ctx context.Context, message string, args ...any) {

}

func (n noopLogger) Warn(ctx context.Context, message string, args ...any) {

}

func (n noopLogger) Error(ctx context.Context, message string, args ...any) {

}
