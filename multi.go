package errors

import "github.com/upfluence/errors/multi"

func WrapErrors(errs []error) error { return WithFrame(multi.Wrap(errs), 1) }
func Combine(errs ...error) error   { return WithFrame(multi.Wrap(errs), 1) }
