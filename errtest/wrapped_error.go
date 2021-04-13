package errtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
)

type ErrorWrapperOptions struct {
	Prefix string
	Suffix string
	N      int

	SkipNil bool
}

func TestErrorWrapper(t *testing.T, fn func(error) error, opts ErrorWrapperOptions) {
	t.Helper()

	for _, tt := range []struct {
		name string
		err  error
		skip bool
	}{
		{name: "nil", skip: opts.SkipNil},
		{name: "sentinel", err: errors.New("scalar")},
		{name: "double wrapped", err: errors.Wrap(errors.New("scalar"), "wrapped")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("skipped by the opts")
			}

			err := fn(tt.err)

			if tt.err == nil {
				assert.Nil(t, err)
				return
			}

			cause := err

			for i := 0; i < opts.N; i++ {
				c, ok := cause.(interface{ Cause() error })
				assert.True(t, ok)
				cause = c.Cause()
			}

			assert.Equal(t, tt.err, cause)

			cause = err

			for i := 0; i < opts.N; i++ {
				u, ok := cause.(interface{ Unwrap() error })
				assert.True(t, ok)
				cause = u.Unwrap()
			}

			assert.Equal(t, tt.err, cause)

			assert.Equal(t, opts.Prefix+tt.err.Error()+opts.Suffix, err.Error())
		})
	}
}
