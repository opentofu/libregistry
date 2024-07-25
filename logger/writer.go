// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package logger

import (
	"context"
	"io"
	"strings"
)

type Level string

const (
	LevelTrace   Level = "TRACE"
	LevelDebug   Level = "DEBUG"
	LevelInfo    Level = "INFO"
	LevelWarning Level = "WARN"
	LevelError   Level = "ERROR"
)

// NewWriter creates an io.WriteCloser that logs each line separately on the given log level, optionally with a prefix.
func NewWriter(ctx context.Context, logger Logger, level Level, prefix string) io.WriteCloser {
	return &writer{
		ctx:    ctx,
		logger: logger,
		level:  level,
		prefix: prefix,
	}
}

type writer struct {
	ctx    context.Context
	logger Logger
	buf    []byte
	level  Level
	prefix string
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)
	lastLineStart := 0
	for i, b := range w.buf {
		if b == 10 || b == 13 {
			line := strings.Trim(string(w.buf[lastLineStart:i]), "\r\n")
			if line != "" {
				w.writeLine(line)
			}
			lastLineStart = i
			// Trim single remaining newline:
			if lastLineStart == len(w.buf)-1 {
				lastLineStart++
			}
		}
	}
	w.buf = w.buf[lastLineStart:]

	return len(p), nil
}

func (w *writer) writeLine(line string) {
	msg := w.prefix + line
	switch w.level {
	case LevelTrace:
		LogTrace(w.ctx, w.logger, "%s", msg)
	case LevelDebug:
		w.logger.Debug(w.ctx, "%s", msg)
	case LevelInfo:
		w.logger.Info(w.ctx, "%s", msg)
	case LevelWarning:
		w.logger.Warn(w.ctx, "%s", msg)
	case LevelError:
		w.logger.Error(w.ctx, "%s", msg)
	default:
		w.logger.Debug(w.ctx, "%s", msg)
	}
}

func (w *writer) Close() error {
	if len(w.buf) > 0 {
		w.writeLine(string(w.buf))
	}
	return nil
}
