package observability

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Exporter is a type of exporter. There will be more exporters
// added in the future.
type Exporter int

const (
	// Stdout is the exporter that prints the trace to stdout.
	Stdout Exporter = iota + 1
	// Jaegar is the exporter that sends the trace to Jaegar.
	Jaegar
)

type endSpanFunc = func(...trace.SpanEndOption)

func StartSpan(ctx context.Context, name string) (context.Context, endSpanFunc) {
	if tracer := GetTracer(ctx); tracer != nil {
		ctx, span := tracer.Start(ctx, name)
		return ctx, span.End
	}
	return ctx, func(...trace.SpanEndOption) {}
}
