package stacktrace_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/stacktrace"
)

func TestGetFramesWithStacktrace(t *testing.T) {
	fs := stacktrace.GetFrames(stacktrace.WithStacktrace(errors.New("foo"), 0, 2))

	assert.Len(t, fs, 3)

	fn, _, line := fs[0].Location()

	assert.Equal(t, "github.com/upfluence/errors/stacktrace_test.TestGetFramesWithStacktrace", fn)
	assert.Equal(t, 13, line)

	fn, _, line = fs[2].Location()

	assert.Equal(t, "github.com/upfluence/errors/stacktrace_test.TestGetFramesWithStacktrace", fn)
	assert.Equal(t, 13, line)
}
