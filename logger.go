package observability

import (
	"context"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

type TeardownFunc func() error

type NewOptions func(context.Context) (context.Context, TeardownFunc)

type zapLoggerKey struct{}

var zapKey zapLoggerKey

// WithZapLogger is an option function to pass a logger
// to the new context.
func WithZapLogger(l *zap.Logger) NewOptions {
	return func(ctx context.Context) (context.Context, TeardownFunc) {
		return context.WithValue(ctx, zapKey, l), nil
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

type tracerKeyT struct{}

var traceKey tracerKeyT

func WithTracing(svcName string, dev bool) NewOptions {
	return func(ctx context.Context) (context.Context, TeardownFunc) {
		// TODO
		if !dev {
			return ctx, nil
		}
		tracer, teardown := newDevTracer(ctx, svcName)
		ctx = context.WithValue(ctx, traceKey, tracer)
		return ctx, func() error {
			return teardown(ctx)
		}
	}
}

func tracer(ctx context.Context) trace.Tracer {
	t, ok := ctx.Value(traceKey).(trace.Tracer)
	if !ok {
		return nil
	}
	return t
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
