package opaque

type opaqueError struct {
	cause error
}

func (oe *opaqueError) Error() string { return oe.cause.Error() }

func (oe *opaqueError) Tags() map[string]interface{} {
	t, ok := oe.cause.(interface{ Tags() map[string]interface{} })

	if !ok {
		return nil
	}

	return t.Tags()
}

func Opaque(err error) error {
	return &opaqueError{cause: err}
}
