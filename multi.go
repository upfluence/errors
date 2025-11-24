package errors

import "github.com/upfluence/errors/multi"

// WrapErrors combines multiple errors into a single error and adds a stack frame.
func WrapErrors(errs []error) error { return WithFrame(multi.Wrap(errs), 1) }

// Combine combines multiple errors into a single error and adds a stack frame.
func Combine(errs ...error) error { return WithFrame(multi.Wrap(errs), 1) }

// Join combines multiple errors into a single error and adds a stack frame.
// This function is compatible with the errors.Join function introduced in Go 1.20.
func Join(errs ...error) error { return WithFrame(multi.Wrap(errs), 1) }
