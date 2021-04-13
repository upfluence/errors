package tags_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/errtest"
	"github.com/upfluence/errors/tags"
)

func TestWithTags(t *testing.T) {
	errtest.TestErrorWrapper(
		t,
		func(err error) error {
			return errors.WithTags(err, map[string]interface{}{"foo": 1})
		},
		errtest.ErrorWrapperOptions{N: 1},
	)
}

func TestGetTags(t *testing.T) {
	assert.Equal(t, 0, len(tags.GetTags(nil)))

	assert.Equal(
		t,
		map[string]interface{}{"domain": "github.com/upfluence/errors/tags", "foo": 1},
		tags.GetTags(
			errors.WithTags(errors.New("foo"), map[string]interface{}{"foo": 1}),
		),
	)
}
