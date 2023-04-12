package observability

import (
	"context"

	"go.uber.org/zap"
)

func NewFromBackground(opts ...NewOptions) (context.Context, teardownFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	for _, apply := range opts {
		ctx = apply(ctx)
	}

	teardown := func() error {
		// TODO: teardown any observability resources
		cancel()
		return nil
	}

	return ctx, teardown
}

func NewFromContext(ctx context.Context, opts ...NewOptions) (context.Context, teardownFunc) {
	for _, apply := range opts {
		ctx = apply(ctx)
	}
	// TODO: teardown any observability resources
	return ctx, func() error { return nil }
}

type teardownFunc func() error

type NewOptions func(context.Context) context.Context

type zapLoggerKey struct{}

var zapKey zapLoggerKey

// WithZapLogger is an option function to pass a logger
// to the new context.
func WithZapLogger(l *zap.Logger) NewOptions {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, zapKey, l)
	}
}

// ZapLogger retrieves the *zap.Logger stored in the context.
// If no logger is on the context, it returns a no-op logger.
// Functionality may be changed in the future to return an error if the
// logger isn't found.
func ZapLogger(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(zapKey).(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}

type fieldsKeyT struct{}

var fieldsKey fieldsKeyT

// GetFields returns all the fields currently set on a given context.
func GetFields(ctx context.Context) []zap.Field {
	if fs, ok := ctx.Value(fieldsKey).([]zap.Field); ok {
		return fs
	}
	return []zap.Field{}
}

// WithFields assigns some fields to a context and returns the new context.
// If the old context has fields already set, any duplicate keys found on the passed map[string]string
// will overwrite the old field values.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	fs := GetFields(ctx)
	for _, field := range fields {
		fs = append(fs, field)
	}

	return context.WithValue(ctx, fieldsKey, fs)
}
