package apperrors

import (
	"errors"
	"fmt"
)

// withStack wraps an error with a captured stack trace.
type withStackError struct {
	cause error
	stack StackTrace
}

func (e *withStackError) Error() string {
	return e.cause.Error()
}

func (e *withStackError) Unwrap() error {
	return e.cause
}

func (e *withStackError) StackTrace() StackTrace {
	return e.stack
}

// Wrap wraps err with a stack trace captured at the call site.
// Returns nil if err is nil.
// If err already contains a *withStack anywhere in its chain, returns err unchanged to prevent double-wrapping.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	var ws *withStackError
	if errors.As(err, &ws) {
		return err
	}
	return &withStackError{cause: err, stack: captureStack(2)}
}

// Wrapf wraps err with a formatted message and a stack trace captured at the call site.
// Returns nil if err is nil.
// If err already contains a *withStack anywhere in its chain, returns err unchanged.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	var ws *withStackError
	if errors.As(err, &ws) {
		return err
	}
	// Append err as the last argument so %w binds to the original error.
	safeArgs := make([]any, len(args)+1)
	copy(safeArgs, args)
	safeArgs[len(args)] = err
	cause := fmt.Errorf(format+": %w", safeArgs...) //nolint:err113 // Wrapf intentionally creates a dynamic error for stack trace enrichment
	return &withStackError{cause: cause, stack: captureStack(2)}
}

// GetStackTrace extracts the StackTrace from err if a *withStack is present in the chain.
// Returns (nil, false) if no stack trace is found.
func GetStackTrace(err error) (StackTrace, bool) {
	var ws *withStackError
	if errors.As(err, &ws) {
		return ws.stack, true
	}
	return nil, false
}
