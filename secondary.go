package errors

import "github.com/upfluence/errors/secondary"

// WithSecondaryError attaches an additional error to the primary error for context.
func WithSecondaryError(err error, additionalErr error) error {
	return secondary.WithSecondaryError(err, additionalErr)
}
