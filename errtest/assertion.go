// Package errtest provides testing utilities for error assertions.
//
// This package offers a fluent interface for asserting error conditions in tests,
// including checking for nil errors, error equality, and error causes.
// It integrates with testify/assert for consistent test failure reporting.
package errtest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
)

// ErrorAssertion represents an assertion to be performed on an error.
type ErrorAssertion interface {
	Assert(testing.TB, error)
}

type multiErrorAssertion []ErrorAssertion

func (eas multiErrorAssertion) Assert(t testing.TB, err error) {
	for _, ea := range eas {
		ea.Assert(t, err)
	}
}

// CombineErrorAssertion combines multiple ErrorAssertion instances into one.
// All assertions will be executed when Assert is called.
func CombineErrorAssertion(eas ...ErrorAssertion) ErrorAssertion {
	if len(eas) == 1 {
		return eas[0]
	}

	return multiErrorAssertion(eas)
}

type noErrorAssertion struct{}

// NoErrorAssertion is an ErrorAssertion that performs no checks.
// Useful as a placeholder or default value.
var NoErrorAssertion ErrorAssertion = noErrorAssertion{}

func (noErrorAssertion) Assert(testing.TB, error) {}

// ErrorAssertionFunc is a function adapter for the ErrorAssertion interface.
type ErrorAssertionFunc func(testing.TB, error)

func (fn ErrorAssertionFunc) Assert(t testing.TB, err error) { fn(t, err) }

// NoError creates an ErrorAssertion that asserts the error is nil.
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

// ErrorEqual creates an ErrorAssertion that asserts the error equals the expected error.
func ErrorEqual(want error, msgAndArgs ...interface{}) ErrorAssertion {
	return ErrorAssertionFunc(
		func(t testing.TB, out error) { assert.Equal(t, want, out, msgAndArgs...) },
	)
}

// ErrorCause creates an ErrorAssertion that asserts the error's root cause
// equals the expected error.
func ErrorCause(want error, msgAndArgs ...interface{}) ErrorAssertion {
	return ErrorAssertionFunc(
		func(t testing.TB, out error) {
			assert.Equal(t, want, errors.Cause(out), msgAndArgs...)
		},
	)
}
