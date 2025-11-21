// Package recovery provides utilities for converting panic recovery values into errors.
//
// This package helps convert the values returned by recover() into proper error types.
// It handles different value types intelligently, preserving runtime.Error instances
// and converting other types into descriptive error messages.
package recovery

import (
	"fmt"
	"runtime"
)

type runtimeError struct {
	v interface{}
}

func (re *runtimeError) RuntimeError() {}

func (re *runtimeError) Error() string {
	switch vv := re.v.(type) {
	case error:
		return vv.Error()
	case fmt.Stringer:
		return vv.String()
	case string:
		return vv
	}

	return fmt.Sprintf("%T: %v", re.v, re.v)
}

// WrapRecoverResult converts a panic recovery value into an error.
// Returns nil if v is nil.
// Preserves runtime.Error instances as-is.
// Converts other values into errors with descriptive messages.
func WrapRecoverResult(v interface{}) error {
	if v == nil {
		return nil
	}

	if rerr, ok := v.(runtime.Error); ok {
		return rerr
	}

	return &runtimeError{v: v}
}
