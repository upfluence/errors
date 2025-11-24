package errors

import "github.com/upfluence/errors/stacktrace"

// WithStack wraps an error with a stack frame captured at the call site.
func WithStack(err error) error {
	return WithFrame(err, 1)
}

// WithStack2 wraps an error with a stack frame captured at the call site, while
// passing through a return value. This is useful for adding stack traces to errors
// from functions that return multiple values (value, error).
//
// If err is nil, returns (v, nil) unchanged. If err is not nil, returns (v, err_with_stack)
// where err_with_stack includes the stack frame.
//
// Example:
//
//	return errors.WithStack2(externalLib.DoSomething())
//
// This is equivalent to but more concise than:
//
//	result, err := externalLib.DoSomething()
//	if err != nil {
//	    return result, errors.WithStack(err)
//	}
//	return result, nil
func WithStack2[T any](v T, err error) (T, error) {
	return v, WithFrame(err, 1)
}

// WithFrame wraps an error with a stack frame at the specified depth in the call stack.
// The depth parameter indicates how many stack frames to skip (0 = current frame).
func WithFrame(err error, d int) error {
	return stacktrace.WithFrame(err, d+1)
}
