// Package stats provides utilities for extracting status strings from errors.
//
// This package allows errors to expose a status string that can be used for
// metrics, logging, or categorization. It traverses the error chain to find
// status information and provides fallback mechanisms when no status is found.
package stats

import (
	"fmt"

	"github.com/upfluence/errors/base"
)

// Statuser provides a custom status string for an error.
type Statuser interface {
	Status(error) string
}

// ExtractStatusOption configures status extraction behavior.
type ExtractStatusOption func(*extractStatusOptions)

type extractStatusOptions struct {
	successStatus    string
	fallbackStatuser Statuser
}

var defaultStatusOptions = extractStatusOptions{
	successStatus:    "success",
	fallbackStatuser: defaultStatuser{},
}

type defaultStatuser struct{}

func (defaultStatuser) Status(err error) string {
	switch t := fmt.Sprintf("%T", err); t {
	case "*errors.errorString", "*errors.fundamental", "*opaque.opaqueError":
		return err.Error()
	default:
		return t
	}
}

// GetStatus extracts a status string from an error.
// Returns the success status string if err is nil.
// Traverses the error chain looking for a Status() method,
// falling back to the configured Statuser if none is found.
func GetStatus(err error, opts ...ExtractStatusOption) string {
	type statuser interface {
		Status() string
	}

	var o = defaultStatusOptions

	for _, opt := range opts {
		opt(&o)
	}

	if err == nil {
		return o.successStatus
	}

	for {
		if st, ok := err.(statuser); ok {
			return st.Status()
		}

		cause := base.UnwrapOnce(err)

		if cause == nil {
			break
		}

		err = cause
	}

	return o.fallbackStatuser.Status(err)
}
