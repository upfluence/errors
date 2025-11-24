// Package secondary provides support for attaching secondary errors to a primary error.
//
// This package is useful when an operation fails and a cleanup or recovery operation
// also fails. The secondary error's tags are included in the combined error, allowing
// additional context to be preserved without losing the primary error information.
package secondary

import (
	"strings"

	"github.com/upfluence/errors/tags"
)

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
	return tags.GetTags(ws.second)
}

func (ws *withSecondary) Errors() []error {
	return []error{ws.cause, ws.second}
}

// WithSecondaryError combines a primary error with a secondary error.
// Returns additionalErr if err is nil, returns err if additionalErr is nil,
// or returns a combined error containing both.
func WithSecondaryError(err error, additionalErr error) error {
	if err == nil {
		return additionalErr
	}

	if additionalErr == nil {
		return err
	}

	return &withSecondary{cause: err, second: additionalErr}
}
