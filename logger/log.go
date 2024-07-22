// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"context"
	"log"
)

// NewGoLogLogger creates a logger which writes to the traditional *log.Logger.
func NewGoLogLogger(logger *log.Logger) Logger {
	return &goLogLogger{
		logger,
	}
}

type goLogLogger struct {
	logger *log.Logger
}

func (g goLogLogger) Trace(_ context.Context, message string, args ...any) {
	g.logger.Printf("TRACE\t"+message, args...)
}

func (g goLogLogger) Debug(_ context.Context, message string, args ...any) {
	g.logger.Printf("DEBUG\t"+message, args...)
}

func (g goLogLogger) Info(_ context.Context, message string, args ...any) {
	g.logger.Printf("INFO\t"+message, args...)
}

func (g goLogLogger) Warn(_ context.Context, message string, args ...any) {
	g.logger.Printf("WARN\t"+message, args...)
}

func (g goLogLogger) Error(_ context.Context, message string, args ...any) {
	g.logger.Printf("ERROR\t"+message, args...)
}
