package errors

import (
	"errors"

	"github.com/upfluence/errors/base"
)

// Cause returns the root cause of an error by recursively unwrapping it.
func Cause(err error) error { return base.UnwrapAll(err) }

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error. Otherwise, Unwrap returns nil.
func Unwrap(err error) error { return base.UnwrapOnce(err) }

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
func As(err error, target interface{}) bool { return errors.As(err, target) }

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool { return errors.Is(err, target) }

type timeout interface {
	Timeout() bool
}

// IsTimeout reports whether any error in err's chain implements the Timeout() bool
// method and returns true from that method.
func IsTimeout(err error) bool {
	var terr timeout

	if As(err, &terr) && terr.Timeout() {
		return true
	}

	return false
}
