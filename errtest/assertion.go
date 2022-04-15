package errtest

import (
	"fmt"
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

type noErrorAssertion struct{}

var NoErrorAssertion ErrorAssertion = noErrorAssertion{}

func (noErrorAssertion) Assert(testing.TB, error) {}

type ErrorAssertionFunc func(testing.TB, error)

func (fn ErrorAssertionFunc) Assert(t testing.TB, err error) { fn(t, err) }

func NoError(msgAndArgs ...interface{}) ErrorAssertion {
	return ErrorAssertionFunc(
		func(t testing.TB, err error) {
			assert.Nil(
				t,
				err,
				append([]interface{}{fmt.Sprint(err)}, msgAndArgs...),
			)
		},
	)
}

func ErrorEqual(want error, msgAndArgs ...interface{}) ErrorAssertion {
	return ErrorAssertionFunc(
		func(t testing.TB, out error) { assert.Equal(t, want, out, msgAndArgs...) },
	)
}

func ErrorCause(want error, msgAndArgs ...interface{}) ErrorAssertion {
	return ErrorAssertionFunc(
		func(t testing.TB, out error) {
			assert.Equal(t, want, errors.Cause(out), msgAndArgs...)
		},
	)
}
