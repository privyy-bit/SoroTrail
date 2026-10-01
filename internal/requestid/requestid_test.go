package requestid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestIDEdgeCases(t *testing.T) {
	// Test empty ID handling
	ctx := context.Background()
	ctxWithEmpty := WithRequestID(ctx, "")
	assert.Equal(t, ctx, ctxWithEmpty, "empty request id should return original context unmodified")
	assert.Equal(t, "", FromContext(ctxWithEmpty))
	assert.Nil(t, Attrs(ctxWithEmpty))

	// Test valid ID setting and retrieval
	validCtx := WithRequestID(ctx, "req-12345")
	assert.Equal(t, "req-12345", FromContext(validCtx))
	attrs := Attrs(validCtx)
	assert.NotNil(t, attrs)
	assert.Len(t, attrs, 2)
	assert.Equal(t, Field, attrs[0])
	assert.Equal(t, "req-12345", attrs[1])

	// Test background job correlation IDs
	jobCtx := WithJob(ctx, JobIngester)
	assert.Equal(t, JobIngester, FromContext(jobCtx))
	jobAttrs := Attrs(jobCtx)
	assert.Contains(t, jobAttrs, JobIngester)

	// Test context cancellation with request id
	cancelCtx, cancel := context.WithCancel(validCtx)
	cancel()
	assert.Equal(t, "req-12345", FromContext(cancelCtx), "cancelled context should still preserve request id")
}
