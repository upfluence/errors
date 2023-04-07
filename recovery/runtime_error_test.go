package recovery

import (
	"errors"
	"net"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapRecoverResultNil(t *testing.T) {
	assert.Nil(t, WrapRecoverResult(nil))
}

func TestWrapRecoverResult(t *testing.T) {
	for _, tt := range []struct {
		in interface{}

		wantError string
	}{
		{in: 1, wantError: "int: 1"},
		{in: errors.New("foo"), wantError: "foo"},
		{in: "foo", wantError: "foo"},
		{in: net.ParseIP("127.0.0.1"), wantError: "127.0.0.1"},
	} {
		err := WrapRecoverResult(tt.in)

		_, ok := err.(runtime.Error)
		assert.True(t, ok)

		assert.Equal(t, tt.wantError, err.Error())
	}
}
