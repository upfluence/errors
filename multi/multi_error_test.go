package multi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/tags"
)

func TestErrorFormatting(t *testing.T) {
	var err = errors.Combine(errors.New("foo"), errors.New("bar"))

	assert.Equal(t, "[foo, bar]", err.Error())
}

func TestTags(t *testing.T) {
	var err = errors.Combine(
		errors.WithTags(errors.New("foo"), map[string]interface{}{"foo": 1}),
		errors.WithTags(errors.New("bar"), map[string]interface{}{"bar": 2}),
	)

	assert.Equal(
		t,
		map[string]interface{}{
			"bar":    2,
			"domain": "github.com/upfluence/errors/multi_test",
			"foo":    1,
		},
		tags.GetTags(err),
	)
}
