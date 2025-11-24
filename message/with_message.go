// Package message provides error wrapping with additional context messages.
//
// This package allows you to wrap errors with descriptive messages while
// preserving the original error. The messages are formatted and prepended
// to the original error message, creating a clear error chain.
package message

import "fmt"

type withMessage struct {
	cause error

	fmt  string
	args []interface{}
}

func (wm *withMessage) Error() string {
	var msg = wm.fmt

	if len(wm.args) > 0 {
		msg = fmt.Sprintf(msg, wm.args...)
	}

	return msg + ": " + wm.cause.Error()
}

func (wm *withMessage) Unwrap() error       { return wm.cause }
func (wm *withMessage) Cause() error        { return wm.cause }
func (wm *withMessage) Args() []interface{} { return wm.args }

// WithMessage wraps an error with an additional context message.
// Returns nil if err is nil.
func WithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &withMessage{cause: err, fmt: msg}
}

// WithMessagef wraps an error with a formatted context message.
// Returns nil if err is nil.
func WithMessagef(err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withMessage{cause: err, fmt: msg, args: args}
}
