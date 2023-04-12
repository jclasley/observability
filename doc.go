// Package observability provides a set of tools to help with observability of applications. It is our experience that most of the meat of a Go service involve a context[context.Context],
// and because of the nature of propagating context, it makes sense to provide an observability API that utilizes a context. Creating an observability context and using it is easy, and convenience methods are
// provided for retrieving information on the context.
//
// The package is designed to be used with [Zap] and [Otel]. There may be more options in the future.
/*
	package main

	func main() {
	  // Create a zap logger
	  logger, _ := zap.NewProduction()

	  // Create a new observability context, using the zap logger option
	  ctx, teardown := observability.NewFromBackground(observability.WithLogger(logger))
	  defer teardown()

	  // Log something
	  observability.Info(ctx, "hello world")

	  // Add some fields to the logger
	  ctx = observability.WithFields(ctx, zap.String("foo", "bar"))

	  // Call something
	  foo(ctx)

	  // Log something else
	  observability.Info(ctx, "goodbye world") // only has {"foo": "bar"}
	}

	func foo(ctx context.Context) {
	  // Add some fields relevant to `foo`
	  ctx = observability.WithFields(ctx, zap.Int("baz", 42)) // has {"foo": "bar", "baz": 42}

	  // Log something
	  observability.Info(ctx, "foo")

	  // ctx can still be used like any ol' context
	  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	  defer cancel()

	  observability.Warn(ctx, "foo is taking a long time", zap.Time("at", time.Now())) // has {"foo": "bar", "baz": 42, "at": "2021-01-01T00:00:00Z"}
	}
*/
// The package also provides some middlewares designed for [Echo].
//
// [Echo]: https://github.com/labstack/echo
// [Zap]: https://go.uber.org/zap
// [Otel]: https://pkg.go.dev/go.opentelemetry.io/otel
package observability
