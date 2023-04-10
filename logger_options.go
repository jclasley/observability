package observability

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

type NewOptions func(context.Context) context.Context

type NewTracingOptions func(*tracingOptions)

type tracingOptions struct {
	exporter Exporter
	resource resource.Resource
}

type zapLoggerKey struct{}

var zapKey zapLoggerKey

// WithZapLogger is an option function to pass a logger
// to the new context.
func WithZapLogger(l *zap.Logger) NewOptions {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, zapKey, l)
	}
}

type tracerKeyT struct{}

var tracerKey tracerKeyT

// WithTracer puts an otel tracer into the context, using the given
// service name as the name of the tracer.
func WithTracer(svcName string, opts ...NewTracingOptions) NewOptions {
	var o tracingOptions
	for _, apply := range opts {
		apply(&o)
	}

	return func(ctx context.Context) context.Context {
		var exp trace.SpanExporter
		var err error
		// TODO: add other exporters
		switch o.exporter {
		case Stdout:
			exp, err = stdouttrace.New()
			if err != nil {
				Error(ctx, fmt.Sprintf("failed to create stdout exporter: %v", err))
			}
		}

		tp := trace.NewTracerProvider(
			trace.WithBatcher(exp),
			trace.WithResource(&o.resource),
		)
		otel.SetTracerProvider(tp)

		tracer := tp.Tracer(svcName)

		return context.WithValue(ctx, tracerKey, tracer)
	}
}

type resourceDefinitionKeyT struct{}

var resourceDefinitionKey resourceDefinitionKeyT

// WithResourceDefinition allows you to specify the resource definition for otel.
// The default resource definition is defined on the ResourceDefinition[ResourceDefinition] type.
func WithResourceDefinition(r resource.Resource) NewTracingOptions {
	return func(opts *tracingOptions) {
		opts.resource = r
	}
}

func defaultResourceDefinition(ctx context.Context) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		if l, ok := ctx.Value(zapKey).(*zap.Logger); ok {
			l.Error("failed to read build info")
		} else {
			log.Printf("failed to read build info")
		}
	}
	// TODO
	fmt.Println(info)
}

// WithExporter allows for controlling which exporter is used. The default is STDOUT.
func WithExporter(e Exporter) NewTracingOptions {
	return func(opts *tracingOptions) {
		opts.exporter = e
	}
}
