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

func Opaque(err error) error {
	return &opaqueError{cause: err}
}
