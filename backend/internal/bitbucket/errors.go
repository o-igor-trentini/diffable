package bitbucket

import (
	"fmt"
	"time"
)

type RateLimitedError struct {
	RetryAfter time.Duration
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("bitbucket: rate limited, retry after %s", e.RetryAfter)
}

type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("bitbucket: %s not found", e.Resource)
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("bitbucket: unauthorized: %s", e.Message)
	}
	return "bitbucket: unauthorized"
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("bitbucket: API error %d: %s", e.StatusCode, e.Message)
}
