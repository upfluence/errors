package errors

import "github.com/upfluence/errors/stats"

// WithStatus attaches a status string to the error and adds a stack frame.
func WithStatus(err error, status string) error {
	return WithFrame(stats.WithStatus(err, status), 1)
}
