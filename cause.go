package errors

import (
	"errors"

	"github.com/upfluence/errors/base"
)

func Cause(err error) error  { return base.UnwrapAll(err) }
func Unwrap(err error) error { return base.UnwrapOnce(err) }

func As(err error, target interface{}) bool { return errors.As(err, target) }
func Is(err, target error) bool             { return errors.Is(err, target) }

type timeout interface {
	Timeout() bool
}

func IsTimeout(err error) bool {
	var terr timeout

	if As(err, &terr) && terr.Timeout() {
		return true
	}

	return false
}
