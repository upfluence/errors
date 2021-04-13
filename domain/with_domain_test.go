package domain_test

import (
	"testing"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/errtest"
)

func TestWithDomain(t *testing.T) {
	errtest.TestErrorWrapper(
		t,
		func(err error) error { return errors.WithDomain(err, "foobar") },
		errtest.ErrorWrapperOptions{N: 1},
	)
}
