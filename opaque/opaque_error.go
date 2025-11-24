// Package opaque provides a way to create opaque errors that hide the underlying
// error type while preserving metadata.
//
// Opaque errors prevent callers from using errors.Is or errors.As to match
// against the wrapped error, while still exposing useful metadata like domain,
// tags, and stacktrace. This is useful for API boundaries where you want to
// expose error information without leaking implementation details.
package opaque

import (
	"github.com/upfluence/errors/domain"
	"github.com/upfluence/errors/stacktrace"
	"github.com/upfluence/errors/tags"
)

type opaqueError struct {
	cause error
}

func (oe *opaqueError) Error() string { return oe.cause.Error() }

func (oe *opaqueError) Domain() domain.Domain {
	return domain.GetDomain(oe.cause)
}

func (oe *opaqueError) Tags() map[string]interface{} {
	return tags.GetTags(oe.cause)
}

func (oe *opaqueError) Frames() []stacktrace.Frame {
	return stacktrace.GetFrames(oe.cause)
}

// Opaque wraps an error to make it opaque, preventing type assertions
// while preserving metadata like domain, tags, and stacktrace.
func Opaque(err error) error {
	return &opaqueError{cause: err}
}
