package domain

import (
	"errors"
	"strings"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("expired token")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

type ValidationErrors []*ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	msgs := make([]string, len(e))
	for i, ve := range e {
		msgs[i] = ve.Error()
	}
	return strings.Join(msgs, "; ")
}
