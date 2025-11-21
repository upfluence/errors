package errors

import "github.com/upfluence/errors/stacktrace"

// WithStack wraps an error with a stack frame captured at the call site.
func WithStack(err error) error {
	return WithFrame(err, 1)
}

// WithFrame wraps an error with a stack frame at the specified depth in the call stack.
// The depth parameter indicates how many stack frames to skip (0 = current frame).
func WithFrame(err error, d int) error {
	return stacktrace.WithFrame(err, d+1)
}
