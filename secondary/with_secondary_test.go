package secondary_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/errtest"
)

func TestWithSecondary(t *testing.T) {
	errtest.TestErrorWrapper(
		t,
		func(err error) error {
			return errors.WithSecondaryError(err, errors.New("bar"))
		},
		errtest.ErrorWrapperOptions{
			N:       1,
			Suffix:  " [ with secondary error: bar]",
			SkipNil: true,
		},
	)

	errw := errors.New("bar")
	assert.Nil(t, errors.WithSecondaryError(nil, nil))
	assert.Equal(t, errw, errors.WithSecondaryError(nil, errw))
}
