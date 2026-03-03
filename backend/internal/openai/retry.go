package openai

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	oai "github.com/sashabaranov/go-openai"
)

const maxRetryAttempts = 3

var baseDelays = []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}

func withRetry(ctx context.Context, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < maxRetryAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if !isRetryable(lastErr) {
			return lastErr
		}

		if attempt < maxRetryAttempts-1 {
			delay := applyJitter(baseDelays[attempt])
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}

	return lastErr
}

func isRetryable(err error) bool {
	apiErr, ok := err.(*oai.APIError)
	if !ok {
		return false
	}
	switch apiErr.HTTPStatusCode {
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusServiceUnavailable:
		return true
	default:
		return false
	}
}

func applyJitter(base time.Duration) time.Duration {
	jitter := 0.8 + rand.Float64()*0.4 // +/- 20%
	return time.Duration(float64(base) * jitter)
}
