package logger

import (
	"context"
)

// Logger is a description of the functions needed for logging in this library.
type Logger interface {
	Trace(ctx context.Context, message string, args ...any)
	Debug(ctx context.Context, message string, args ...any)
	Info(ctx context.Context, message string, args ...any)
	Warn(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, message string, args ...any)
}
