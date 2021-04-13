package inhibit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors/reporter"
)

type mockReporter struct {
	called bool
}

func (*mockReporter) Close() error { return nil }
func (mr *mockReporter) Report(error, reporter.ReportOptions) {
	mr.called = true
}

func TestReporter(t *testing.T) {
	var err1 = errors.New("foo")

	for _, tt := range []struct {
		err    error
		called bool
	}{
		{err: err1},
		{err: errors.New("bar"), called: true},
	} {
		var (
			mr mockReporter

			r = NewReporter(&mr)
		)

		r.AddErrorInhibitors(
			ErrorInhibitorFunc(func(err error) bool { return err == err1 }),
		)

		r.Report(tt.err, reporter.ReportOptions{})

		assert.Equal(t, tt.called, mr.called)
	}

}
