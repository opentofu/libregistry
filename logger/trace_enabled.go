//go:build log_trace

package logger

// LogTrace logs a trace message using the specified logger. This function should be used where the trace logs
// may need to be disabled using a build flag. The purpose of this function is to enable inlining the noop function if
// the build tag is disabled.
func LogTrace(ctx context.Context, logger Logger, message string, args ...any) {
	logger.Trace(ctx, message, args...)
}
