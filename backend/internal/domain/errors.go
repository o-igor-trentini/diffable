package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound           = errors.New("resource not found")
	ErrValidation         = errors.New("validation error")
	ErrExternalService    = errors.New("external service error")
	ErrRateLimited        = errors.New("rate limit exceeded")
	ErrTimeout            = errors.New("request timeout")
	ErrTokenLimitExceeded = errors.New("token limit exceeded")
)

// WrapError wraps a sentinel error with additional context.
func WrapError(sentinel error, msg string) error {
	return fmt.Errorf("%s: %w", msg, sentinel)
}
