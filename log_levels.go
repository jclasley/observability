package observability

import (
	"context"

	"go.uber.org/zap"
)

// getLoggerAssignFields takes a context and attempts to lookup the logger and fields,
// then sets the fields on a new logger and returns that, ready for logging
func getLoggerAssignFields(ctx context.Context, fields ...zap.Field) *zap.Logger {
	l := GetZapLogger(ctx)

	fs := Fields(ctx)
	for _, field := range fields {
		fs = append(fs, field)
	}
	return l.With(fs...)
}

// Info writes an info level log message with the msg. Optional fields
// can be included for use in only this message
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	l := getLoggerAssignFields(ctx, fields...)
	l.Info(msg)
}

// Debug writes a debug level log message with the msg. Optional fields
// can be included for use in only this message.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l := getLoggerAssignFields(ctx, fields...)
	l.Debug(msg)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l := getLoggerAssignFields(ctx, fields...)
	l.Warn(msg)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	l := getLoggerAssignFields(ctx, fields...)
	l.Error(msg)
}
