package domain

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrValidation      = errors.New("validation error")
	ErrExternalService = errors.New("external service error")
)
