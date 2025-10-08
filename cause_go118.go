//go:build go1.18

package errors

import "errors"

func IsOfType[T error](err error) bool {
	for {
		if _, ok := err.(T); ok || err == nil {
			return ok
		}

		err = Unwrap(err)
	}
}

func AsType[T error](err error) (T, bool) {
	var e T

	return e, errors.As(err, &e)
}
