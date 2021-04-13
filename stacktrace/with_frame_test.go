package stacktrace_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/errtest"
	"github.com/upfluence/errors/stacktrace"
)

func TestWithFrame(t *testing.T) {
	errtest.TestErrorWrapper(
		t,
		func(err error) error { return errors.WithStack(err) },
		errtest.ErrorWrapperOptions{N: 1},
	)
}

func TestGetFrames(t *testing.T) {
	fs := stacktrace.GetFrames(errors.New("foo"))

	assert.Len(t, fs, 1)

	fn, _, line := fs[0].Location()

	assert.Equal(t, "github.com/upfluence/errors/stacktrace_test.TestGetFrames", fn)
	assert.Equal(t, 22, line)
}
