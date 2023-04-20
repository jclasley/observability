package observability

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Middlewares(logger *zap.Logger, options ...LoggingOpts) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		LoggingMiddleware(logger, options...),
	}
}

func LoggingMiddleware(logger *zap.Logger, options ...LoggingOpts) echo.MiddlewareFunc {
	var opts LoggingOptions
	for _, apply := range options {
		apply(&opts)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, teardown := NewFromContext(c.Request().Context(), WithZapLogger(logger))
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

			// call other middlewares
			err := next(c)

			fields = append(fields, zap.Int("status", c.Response().Status), zap.String("method", c.Request().Method), zap.Duration("dur", time.Since(now)))

			if err != nil || c.Response().Status >= 400 {
				Error(ctx, c.Path(), fields...)
			} else {
				Info(ctx, c.Path(), fields...)
			}
			return err
		}
	}
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
