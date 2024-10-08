// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package logger

import (
	"context"
)

// Logger is a description of the functions needed for logging in this library.
type Logger interface {
	// WithName creates a new named logger.
	WithName(name string) Logger
	Trace(ctx context.Context, message string, args ...any)
	Debug(ctx context.Context, message string, args ...any)
	Info(ctx context.Context, message string, args ...any)
	Warn(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, message string, args ...any)
}
