package observability

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Middlewares(ctx context.Context, logger *zap.Logger, options ...LoggingOpts) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		LoggingMiddleware(ctx, logger, options...),
	}
}

func LoggingMiddleware(ctx context.Context, logger *zap.Logger, options ...LoggingOpts) echo.MiddlewareFunc {
	var opts LoggingOptions
	for _, apply := range options {
		apply(&opts)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, teardown := NewFromContext(ctx, WithZapLogger(logger))
			defer teardown()

			now := time.Now()

			fields := []zap.Field{}
			if opts.RequestID {
				id, err := uuid.NewUUID()
				if err != nil {
					Warn(ctx, "failed to generate request id", zap.Error(err))
				} else {
					ctx = WithFields(ctx, zap.String("req_id", id.String()))
				}
			}

			req := c.Request()
			c.SetRequest(req.WithContext(ctx))

			if opts.Timestamp {
				fields = append(fields, zap.Time("at", time.Now()))
			}

			rw := newWriter(c.Response())
			c.Response().Writer = rw

			// call other middlewares
			err := next(c)

			fields = append(fields, zap.Int("status", rw.code), zap.String("method", c.Request().Method), zap.Duration("dur", time.Since(now)))

			if err != nil || rw.code >= 400 {
				Error(ctx, c.Path(), fields...)
			} else {
				Info(ctx, c.Path(), fields...)
			}
			return err
		}
	}
}

type responseWriter struct {
	code int
	rw   echo.Response
}

func newWriter(rw *echo.Response) *responseWriter {
	return &responseWriter{rw: *rw}
}

func (w *responseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.rw.Write(b)
}

func (w *responseWriter) WriteHeader(code int) {
	w.code = code
	w.rw.WriteHeader(code)
}

type LoggingOptions struct {
	Timestamp bool
	RequestID bool
}

type LoggingOpts func(opts *LoggingOptions)

func WithTimestamp() LoggingOpts {
	return func(opts *LoggingOptions) {
		opts.Timestamp = true
	}
}

func WithRequestID() LoggingOpts {
	return func(opts *LoggingOptions) {
		opts.RequestID = true
	}
}
