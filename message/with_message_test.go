package message_test

import (
	"testing"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/errtest"
)

func TestWrap(t *testing.T) {
	errtest.TestErrorWrapper(
		t,
		func(err error) error { return errors.Wrap(err, "foo") },
		errtest.ErrorWrapperOptions{Prefix: "foo: ", N: 2},
	)
}

func TestWrapf(t *testing.T) {
	errtest.TestErrorWrapper(
		t,
		func(err error) error { return errors.Wrapf(err, "foo %q", "bar") },
		errtest.ErrorWrapperOptions{Prefix: "foo \"bar\": ", N: 2},
	)
}
