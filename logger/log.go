// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package logger

import (
	"context"
	"log"
)

// NewGoLogLogger creates a logger which writes to the traditional *log.Logger.
func NewGoLogLogger(logger *log.Logger) Logger {
	return &goLogLogger{
		"",
		logger,
	}
}

type goLogLogger struct {
	prefix string
	logger *log.Logger
}

func (g goLogLogger) WithName(name string) Logger {
	return &goLogLogger{
		name + "\t",
		g.logger,
	}
}

func (g goLogLogger) log(_ context.Context, level string, message string, args ...any) {
	if g.prefix != "" {
		g.logger.Printf(g.prefix+level+"\t"+message, args...)
	} else {
		g.logger.Printf(level+"\t"+message, args...)
	}
}

func (g goLogLogger) Trace(ctx context.Context, message string, args ...any) {
	g.log(ctx, "TRACE", message, args...)
}

func (g goLogLogger) Debug(ctx context.Context, message string, args ...any) {
	g.log(ctx, "DEBUG", message, args...)
}

func (g goLogLogger) Info(ctx context.Context, message string, args ...any) {
	g.log(ctx, "INFO", message, args...)
}

func (g goLogLogger) Warn(ctx context.Context, message string, args ...any) {
	g.log(ctx, "WARN", message, args...)
}

func (g goLogLogger) Error(ctx context.Context, message string, args ...any) {
	g.log(ctx, "ERROR", message, args...)
}
