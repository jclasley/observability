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

	ctx := NewFromBackground(WithZapLogger(l))
	ctxL := ZapLogger(ctx)

	require.Same(t, l, ctxL)
}

func TestFields(t *testing.T) {
	ctx := WithFields(context.Background(), map[string]string{"some": "value", "goes": "here"})
	ctx = WithFields(ctx, map[string]string{"more": "fields"})

	require.Equal(t, Fields(ctx), map[string]string{"some": "value", "goes": "here", "more": "fields"})
}

