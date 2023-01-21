//go:build go1.18

package errors

func IsOfType[T error](err error) bool {
	for {
		if _, ok := err.(T); ok || err == nil {
			return ok
		}

		err = Unwrap(err)
	}
}
