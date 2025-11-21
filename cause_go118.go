//go:build go1.18

package errors

import "errors"

// IsOfType reports whether any error in err's tree matches the generic type T.
// It traverses the error chain using Unwrap until it finds a match or reaches the end.
func IsOfType[T error](err error) bool {
	for {
		if _, ok := err.(T); ok || err == nil {
			return ok
		}

		err = Unwrap(err)
	}
}

// AsType attempts to convert err to the generic type T by traversing the error chain.
// It returns the converted error and true if successful, or a zero value and false otherwise.
// This is a type-safe wrapper around errors.As.
func AsType[T error](err error) (T, bool) {
	var e T

	return e, errors.As(err, &e)
}
