package logger

import (
	"context"
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

func (s slogLogger) Trace(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, -8, message, args...)
}

func (s slogLogger) Debug(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, slog.LevelDebug, message, args...)
}

func (s slogLogger) Info(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, slog.LevelInfo, message, args...)
}

func (s slogLogger) Warn(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, slog.LevelWarn, message, args...)
}

func (s slogLogger) Error(ctx context.Context, message string, args ...any) {
	s.backingLogger.Log(ctx, slog.LevelError, message, args...)
}
