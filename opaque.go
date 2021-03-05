package errors

import "github.com/upfluence/errors/opaque"

func Opaque(err error) error { return opaque.Opaque(err) }
