package secondary

import "strings"

type withSecondary struct {
	cause  error
	second error
}

func (ws *withSecondary) Error() string {
	var b strings.Builder

	b.WriteString(ws.cause.Error())
	b.WriteString(" [ with secondary error: ")
	b.WriteString(ws.second.Error())
	b.WriteString("]")

	return b.String()
}

func (ws *withSecondary) Unwrap() error { return ws.cause }
func (ws *withSecondary) Cause() error  { return ws.cause }

func (ws *withSecondary) Tags() map[string]interface{} {
	if t, ok := ws.second.(interface{ Tags() map[string]interface{} }); ok {
		return t.Tags()
	}

	return nil
}

func (ws *withSecondary) Errors() []error {
	return []error{ws.cause, ws.second}
}

func WithSecondaryError(err error, additionalErr error) error {
	if err == nil {
		return additionalErr
	}

	if additionalErr == nil {
		return err
	}

	return &withSecondary{cause: err, second: additionalErr}
}
