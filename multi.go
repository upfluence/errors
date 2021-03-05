package errors

import "github.com/upfluence/errors/multi"

func WrapErrors(errs []error) error { return multi.Wrap(errs) }
func Combine(errs ...error) error   { return multi.Wrap(errs) }
