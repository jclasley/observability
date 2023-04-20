package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSameLogger(t *testing.T) {
	l, err := zap.NewDevelopment()
	require.NoError(t, err)

	ctx, teardown := NewFromBackground(WithZapLogger(l))
	defer teardown()

	ctxL := ZapLogger(ctx)

	require.Same(t, l, ctxL)
}

func TestNoLogger(t *testing.T) {
	l := ZapLogger(context.Background())
	require.NotNil(t, l)
}

func TestFields(t *testing.T) {
	ctx := WithFields(context.Background(), zap.String("something", "here"), zap.Int("x", 1))
	ctx = WithFields(ctx, zap.Bool("false", false))

	require.Equal(t, GetFields(ctx), []zap.Field{zap.String("something", "here"), zap.Int("x", 1), zap.Bool("false", false)})
}
