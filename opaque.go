package errors

import "github.com/upfluence/errors/opaque"

// Opaque wraps an error to prevent unwrapping, hiding the underlying error chain.
func Opaque(err error) error { return opaque.Opaque(err) }
