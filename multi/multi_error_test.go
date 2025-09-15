package multi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/multi"
	"github.com/upfluence/errors/tags"
)

func TestErrorFormatting(t *testing.T) {
	var err = errors.Combine(errors.New("foo"), errors.New("bar"))

	assert.Equal(t, "[foo, bar]", err.Error())
}

func TestCombining(t *testing.T) {
	foo := errors.New("foo")

	for _, tt := range []struct {
		errs []error

		want error
	}{
		{errs: []error{nil, nil, errors.Combine(nil, nil)}},
		{errs: []error{}},
		{errs: []error{nil, nil, errors.Combine(nil, foo)}, want: foo},
		{errs: []error{foo}, want: foo},
		{
			errs: []error{
				foo, errors.Combine(foo, foo),
			},
			want: multi.Combine(foo, foo, foo),
		},
		{
			errs: []error{
				errors.WithStack(notComparableError{s: []string{"foo"}}),
				errors.Combine(
					notComparableError{s: []string{"bar"}},
					notComparableError{s: []string{"buz"}},
				),
			},
			want: multi.Combine(
				notComparableError{s: []string{"foo"}},
				notComparableError{s: []string{"bar"}},
				notComparableError{s: []string{"buz"}},
			),
		},
	} {
		assert.Equal(t, tt.want, errors.Cause(errors.WrapErrors(tt.errs)))
	}
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

type notComparableError struct {
	s []string
}

func (notComparableError) Error() string { return "not comparable" }
