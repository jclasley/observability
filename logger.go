package observability

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TODO: add in a teardown func
// NewFromBackground creates a new context from the background context,
// using passed in options to add things to the context for observability
// purposes.
func NewFromBackground(opts ...NewOptions) context.Context {
	ctx := context.Background()
	for _, apply := range opts {
		ctx = apply(ctx)
	}

	if _, ok := ctx.Value(resourceDefinitionKey).(resource.Resource); !ok {
		// get default resource definitions
	}

	return ctx
}

// NewFromContext creates a new context from the given context,
// using passed in options to add things to the context for observability
// purposes.
func NewFromContext(ctx context.Context, opts ...NewOptions) context.Context {
	for _, apply := range opts {
		ctx = apply(ctx)
	}
	return ctx
}

// GetZapLogger retrieves the *zap.Logger stored in the context.
// If no logger is on the context, it returns a global logger.
// Functionality may be changed in the future to return an error if the
// logger isn't found.
func GetZapLogger(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(zapKey).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}

// GetTracer returns the tracer stored in the context.
// If not tracer is on the context, it returns nil.
// This may someday return a no-op tracer.
func GetTracer(ctx context.Context) trace.Tracer {
	if tracer, ok := ctx.Value(tracerKey).(trace.Tracer); ok {
		return tracer
	}
	return nil
}

type fieldsKeyT struct{}

var fieldsKey fieldsKeyT

// Fields returns all the fields currently set on a given context.
func Fields(ctx context.Context) []zap.Field {
	if fs, ok := ctx.Value(fieldsKey).([]zap.Field); ok {
		return fs
	}
	return []zap.Field{}
}

// WithFields assigns some fields to a context and returns the new context.
// If the old context has fields already set, any duplicate keys found on the passed map[string]string
// will overwrite the old field values.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	fs := Fields(ctx)
	for _, field := range fields {
		fs = append(fs, field)
	}

	return context.WithValue(ctx, fieldsKey, fs)
}
