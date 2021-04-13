package stats_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/errtest"
	"github.com/upfluence/errors/stats"
	"github.com/upfluence/errors/tags"
)

func TestWithStauts(t *testing.T) {
	errtest.TestErrorWrapper(
		t,
		func(err error) error {
			return errors.WithStatus(err, "foo")
		},
		errtest.ErrorWrapperOptions{N: 1},
	)
}

func TestGetStatus(t *testing.T) {
	assert.Equal(
		t,
		"success",
		stats.GetStatus(nil),
	)

	assert.Equal(
		t,
		"baz",
		stats.GetStatus(errors.New("baz")),
	)

	assert.Equal(
		t,
		"foo",
		stats.GetStatus(errors.WithStatus(errors.New("bar"), "foo")),
	)
}

func TestGetTags(t *testing.T) {
	assert.Equal(
		t,
		map[string]interface{}{
			"domain": "github.com/upfluence/errors/stats_test",
			"status": "baz",
		},
		tags.GetTags(stats.WithStatus(errors.New("bar"), "baz")),
	)
}
