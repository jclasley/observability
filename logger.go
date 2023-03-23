package observability

import (
	"context"
	"go.uber.org/zap"
)

func NewFromBackground(opts ...NewOptions) context.Context {
	ctx := context.Background()
	for _, apply := range opts {
		ctx = apply(ctx)
	}
	return ctx
}

type NewOptions func (context.Context) context.Context

type zapLoggerKey struct{}
var zapKey struct{}

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

// Fields returns all the fields currently set on a given context.
func Fields(ctx context.Context) map[string]string {
	if fs, ok := ctx.Value(fieldsKey).(map[string]string); ok {
		return fs
	}
	return map[string]string{}
}

// WithFields assigns some fields to a context and returns the new context.
// If the old context has fields already set, any duplicate keys found on the passed map[string]string
// will overwrite the old field values.
func WithFields(ctx context.Context, fields map[string]string) context.Context {
	fs := Fields(ctx)
	for k, v := range fields {
		fs[k] = v
	}
	return context.WithValue(ctx, fieldsKey, fs)
}

func Log(ctx context.Context, msg string, fields ...map[string]string) {
	l := ZapLogger(ctx)
	fs := Fields(ctx)

	ctxZapFields := convertToZapField(fs)
	zapFields := make([]zap.Field, 0)
	for _, v := range fields {
		zapFields = append(zapFields, convertToZapField(v)...)
	}
	ctxZapFields = append(ctxZapFields, zapFields...)

	l.Info(msg, ctxZapFields...)
}

func convertToZapField(m map[string]string) []zap.Field {
	fs := make([]zap.Field, 0, len(m))
	for k, v := range m {
		fs = append(fs, zap.String(k, v))
	}
	return fs
}
