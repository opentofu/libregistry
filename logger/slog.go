// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package logger

import (
	"context"
	"fmt"
	"log/slog"
)

// NewSLogLogger creates a logger that writes to an *slog.Logger backend.
func NewSLogLogger(backingLogger *slog.Logger) Logger {
	return &slogLogger{
		backingLogger,
	}
}

type slogLogger struct {
	backingLogger *slog.Logger
}

func (s slogLogger) WithName(name string) Logger {
	return &slogLogger{
		s.backingLogger.WithGroup(name),
	}
}

func (s slogLogger) Trace(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, -8, fmt.Sprintf(message, args...))
}

func (s slogLogger) Debug(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, slog.LevelDebug, fmt.Sprintf(message, args...))
}

func (s slogLogger) Info(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, slog.LevelInfo, fmt.Sprintf(message, args...))
}

func (s slogLogger) Warn(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, slog.LevelWarn, fmt.Sprintf(message, args...))
}

func (s slogLogger) Error(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, slog.LevelError, fmt.Sprintf(message, args...))
}
