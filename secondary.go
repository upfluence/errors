package errors

import "github.com/upfluence/errors/secondary"

func WithSecondaryError(err error, additionalErr error) error {
	return secondary.WithSecondaryError(err, additionalErr)
}
