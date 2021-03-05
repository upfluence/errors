package errtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
)

type ErrorAssertion interface {
	Assert(testing.TB, error)
}

type multiErrorAssertion []ErrorAssertion

func (eas multiErrorAssertion) Assert(t testing.TB, err error) {
	for _, ea := range eas {
		ea.Assert(t, err)
	}
}

func CombineErrorAssertion(eas ...ErrorAssertion) ErrorAssertion {
	if len(eas) == 1 {
		return eas[0]
	}

	return multiErrorAssertion(eas)
}

type NoErrorAssertion struct{}

func NoError(msgAndArgs ...interface{}) ErrorAssertion {
	return func(t testing.TB, err error) { assert.Nil(t, err, msgAndArgs...) }
}

func ErrorEqual(want error, msgAndArgs ...interface{}) ErrorAssertion {
	return func(t testing.TB, out error) {
		assert.Equal(t, want, out, msgAndArgs...)
	}
}

func ErrorCause(want error, msgAndArgs ...interface{}) ErrorAssertion {
	return func(t testing.TB, out error) {
		assert.Equal(t, want, errors.Cause(out), msgAndArgs...)
	}
}
