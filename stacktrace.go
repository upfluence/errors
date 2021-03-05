package errors

import "github.com/upfluence/errors/stacktrace"

func WithStack(err error) error {
	return WithFrame(err, 1)
}

func WithFrame(err error, d int) error {
	return stacktrace.WithFrame(err, d+1)
}
