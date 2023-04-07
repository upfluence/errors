package errors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTimeout bool

func (mockTimeout) Error() string    { return "mock" }
func (mt mockTimeout) Timeout() bool { return bool(mt) }

func TestIsTimeout(t *testing.T) {
	for _, tt := range []struct {
		input error

		isTimeout bool
	}{
		{},
		{
			input:     New("foo"),
			isTimeout: false,
		},
		{
			input:     Wrap(New("foo"), "foo"),
			isTimeout: false,
		},
		{
			input:     context.DeadlineExceeded,
			isTimeout: true,
		},
		{
			input:     Wrap(context.DeadlineExceeded, "wrapped"),
			isTimeout: true,
		},
		{
			input:     Opaque(context.DeadlineExceeded),
			isTimeout: false,
		},
		{
			input:     mockTimeout(true),
			isTimeout: true,
		},
		{
			input:     Wrap(mockTimeout(true), "wrapped"),
			isTimeout: true,
		},
		{
			input:     mockTimeout(false),
			isTimeout: false,
		},
	} {
		assert.Equal(t, tt.isTimeout, IsTimeout(tt.input))
	}
}
