package logger

import (
	"context"
	"testing"
)

// NewTestLogger produces a new logger that writes to a Go *testing.T.
func NewTestLogger(t *testing.T) Logger {
	return &testLogger{
		t,
	}
}

type testLogger struct {
	t *testing.T
}

func (t testLogger) Trace(_ context.Context, message string, args ...any) {
	t.t.Logf("TRACE\t"+message, args...)
}

func (t testLogger) Debug(_ context.Context, message string, args ...any) {
	t.t.Logf("DEBUG\t"+message, args...)
}

func (t testLogger) Info(_ context.Context, message string, args ...any) {
	t.t.Logf("INFO\t"+message, args...)
}

func (t testLogger) Warn(_ context.Context, message string, args ...any) {
	t.t.Logf("WARN\t"+message, args...)
}

func (t testLogger) Error(_ context.Context, message string, args ...any) {
	t.t.Logf("ERROR\t"+message, args...)
}
