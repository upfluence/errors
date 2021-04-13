package opaque

import "github.com/upfluence/errors/tags"

type opaqueError struct {
	cause error
}

func (oe *opaqueError) Error() string { return oe.cause.Error() }

func (oe *opaqueError) Tags() map[string]interface{} {
	return tags.GetTags(oe.cause)
}

func Opaque(err error) error {
	return &opaqueError{cause: err}
}
