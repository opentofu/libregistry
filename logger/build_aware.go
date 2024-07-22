// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"context"
)

// NewBuildAwareLogger returns a logger that is able to skip trace logs when built without the log_trace build tag.
func NewBuildAwareLogger(backingLogger Logger) Logger {
	return &buildAwareLogger{
		backingLogger,
	}
}

type buildAwareLogger struct {
	backingLogger Logger
}

func (b buildAwareLogger) WithName(name string) Logger {
	return &buildAwareLogger{
		b.backingLogger.WithName(name),
	}
}

func (b buildAwareLogger) Trace(ctx context.Context, message string, args ...any) {
	LogTrace(ctx, b.backingLogger, message, args...)
}

func (b buildAwareLogger) Debug(ctx context.Context, message string, args ...any) {
	b.backingLogger.Debug(ctx, message, args)
}

func (b buildAwareLogger) Info(ctx context.Context, message string, args ...any) {
	b.backingLogger.Info(ctx, message, args)
}

func (b buildAwareLogger) Warn(ctx context.Context, message string, args ...any) {
	b.backingLogger.Warn(ctx, message, args)
}

func (b buildAwareLogger) Error(ctx context.Context, message string, args ...any) {
	b.backingLogger.Error(ctx, message, args)
}
