package observability

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	otrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

type teardown = func(context.Context) error

func newDevTracer(ctx context.Context, svcName string) (trace.Tracer, teardown) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			semconv.ServiceName(svcName),
		),
	)
	if err != nil {
		panic(err)
	}

	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		panic(err)
	}

	tp := otrace.NewTracerProvider(otrace.WithBatcher(exporter), otrace.WithResource(r))
	tracer := tp.Tracer(svcName)
	return tracer, tp.Shutdown
}

type Attributes map[string]string

type spanKeyT struct{}

var spanKey spanKeyT

func StartSpan(ctx context.Context, name string, attrs ...Attributes) (context.Context, func(...trace.SpanEndOption)) {
	t := tracer(ctx)
	if t == nil {
		return ctx, func(_ ...trace.SpanEndOption) {}
	}

	traceAttr := mapToKeyValue(attrs...)

	ctx, span := t.Start(ctx, name, trace.WithAttributes(traceAttr...))
	ctx = context.WithValue(ctx, spanKey, span)
	return ctx, span.End
}

func mapToKeyValue(attrs ...Attributes) []attribute.KeyValue {
	concatAttr := make(map[string]string)
	for _, m := range attrs {
		for k, v := range m {
			concatAttr[k] = v
		}
	}

	traceAttr := make([]attribute.KeyValue, 0, len(concatAttr))
	for k, v := range concatAttr {
		traceAttr = append(traceAttr, attribute.String(k, v))
	}
	return traceAttr
}

func ctxSpan(ctx context.Context) trace.Span {
	span, ok := ctx.Value(spanKey).(trace.Span)
	if !ok {
		return nil
	}
	return span
}

type attrKeyT struct{}

var attrKey attrKeyT

func WithAttributes(ctx context.Context, attrs ...Attributes) context.Context {
	ctxAttrs := ctxAttributes(ctx)
	attrs = append(ctxAttrs, attrs...)

	span := ctxSpan(ctx)
	if span != nil {
		traceAttr := mapToKeyValue(attrs...)
		span.SetAttributes(traceAttr...)
	}

	ctx = context.WithValue(ctx, attrKey, attrs)
	return ctx
}

func ctxAttributes(ctx context.Context) []Attributes {
	a, ok := ctx.Value(attrKey).([]Attributes)
	if !ok {
		return nil
	}
	return a
}
