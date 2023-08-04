package observability

import (
	"context"
	"errors"
)

func NewFromBackground(opts ...NewOptions) (context.Context, TeardownFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	teardowns := make([]TeardownFunc, len(opts))
	for k, apply := range opts {
		var td TeardownFunc
		ctx, td = apply(ctx)
		teardowns[k] = td
	}

	teardown := func() error {
		defer cancel()

		var err error
		for _, td := range teardowns {
			innerErr := td()
			err = errors.Join(err, innerErr)
		}
		return err
	}

	return ctx, teardown
}

func NewFromContext(ctx context.Context, opts ...NewOptions) (context.Context, TeardownFunc) {
	ctx, cancel := context.WithCancel(ctx)

	teardowns := make([]TeardownFunc, len(opts))
	for k, apply := range opts {
		var td TeardownFunc
		ctx, td = apply(ctx)
		teardowns[k] = td
	}

	teardown := func() error {
		defer cancel()

		var err error
		for _, td := range teardowns {
			innerErr := td()
			err = errors.Join(err, innerErr)
		}
		return err
	}
	return ctx, teardown
}
