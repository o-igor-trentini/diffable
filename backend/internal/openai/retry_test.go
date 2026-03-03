package openai

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	oai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithRetry_SuccessFirstAttempt(t *testing.T) {
	calls := 0
	err := withRetry(context.Background(), func() error {
		calls++
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}

func TestWithRetry_SuccessSecondAttempt(t *testing.T) {
	calls := 0
	start := time.Now()

	err := withRetry(context.Background(), func() error {
		calls++
		if calls == 1 {
			return &oai.APIError{HTTPStatusCode: http.StatusTooManyRequests, Message: "rate limited"}
		}
		return nil
	})

	elapsed := time.Since(start)
	require.NoError(t, err)
	assert.Equal(t, 2, calls)
	assert.GreaterOrEqual(t, elapsed, 700*time.Millisecond) // ~1s with jitter
}

func TestWithRetry_FailsAfterMaxAttempts(t *testing.T) {
	calls := 0

	err := withRetry(context.Background(), func() error {
		calls++
		return &oai.APIError{HTTPStatusCode: http.StatusInternalServerError, Message: "server error"}
	})

	require.Error(t, err)
	assert.Equal(t, 3, calls)
	var apiErr *oai.APIError
	assert.ErrorAs(t, err, &apiErr)
}

func TestWithRetry_NonRetryableError(t *testing.T) {
	calls := 0

	err := withRetry(context.Background(), func() error {
		calls++
		return &oai.APIError{HTTPStatusCode: http.StatusBadRequest, Message: "bad request"}
	})

	require.Error(t, err)
	assert.Equal(t, 1, calls) // Does not retry on 400
}

func TestWithRetry_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := withRetry(ctx, func() error {
		return &oai.APIError{HTTPStatusCode: http.StatusTooManyRequests}
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
}

func TestWithRetry_503Retryable(t *testing.T) {
	calls := 0

	err := withRetry(context.Background(), func() error {
		calls++
		if calls < 3 {
			return &oai.APIError{HTTPStatusCode: http.StatusServiceUnavailable}
		}
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 3, calls)
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"429 is retryable", &oai.APIError{HTTPStatusCode: 429}, true},
		{"500 is retryable", &oai.APIError{HTTPStatusCode: 500}, true},
		{"503 is retryable", &oai.APIError{HTTPStatusCode: 503}, true},
		{"400 is not retryable", &oai.APIError{HTTPStatusCode: 400}, false},
		{"401 is not retryable", &oai.APIError{HTTPStatusCode: 401}, false},
		{"404 is not retryable", &oai.APIError{HTTPStatusCode: 404}, false},
		{"generic error is not retryable", errors.New("generic"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isRetryable(tt.err))
		})
	}
}
