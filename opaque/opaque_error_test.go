package opaque_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/tags"
)

func TestOpaque(t *testing.T) {
	errw := errors.WithTags(errors.New("bar"), map[string]interface{}{"foo": 1})
	err := errors.Opaque(errw)

	assert.Equal(t, err, errors.Cause(err))
	assert.Equal(
		t,
		map[string]interface{}{"domain": "github.com/upfluence/errors/opaque", "foo": 1},
		tags.GetTags(err),
	)
	assert.Equal(t, "bar", err.Error())
}
