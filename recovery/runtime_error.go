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

func WrapRecoverResult(v interface{}) error {
	if v == nil {
		return nil
	}

	if rerr, ok := v.(runtime.Error); ok {
		return rerr
	}

	return &runtimeError{v: v}
}
