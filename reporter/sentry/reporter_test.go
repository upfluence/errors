package sentry

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/reporter"
)

func TestBuildEvent(t *testing.T) {
	for _, tt := range []struct {
		name string

		opts  []Option
		err   error
		ropts reporter.ReportOptions

		evtfn func(*testing.T, *sentry.Event)
	}{
		{
			name: "no error",
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assert.Nil(t, evt)
			},
		},
		{
			name: "simple error",
			err:  errors.New("basic error"),
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assert.Equal(t, "basic error", evt.Message)
				assert.Equal(t, "", evt.Transaction)
				assert.Equal(
					t,
					map[string]string{"domain": "github.com/upfluence/errors/reporter/sentry"},
					evt.Tags,
				)
				assert.Equal(t, map[string]interface{}{}, evt.Contexts)
			},
			ropts: reporter.ReportOptions{Depth: 1},
		},
		{
			name: "error transaction tag",
			err: errors.WithTags(
				errors.New("basic error"),
				map[string]interface{}{reporter.TransactionKey: "transaction#27"},
			),
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assert.Equal(t, "basic error", evt.Message)
				assert.Equal(t, "transaction#27", evt.Transaction)
			},
		},
		{
			name: "simple error with thrift opts",
			err:  errors.New("basic error"),
			ropts: reporter.ReportOptions{
				Tags: map[string]interface{}{
					reporter.ThriftRequestServiceKey: "svc",
					reporter.ThriftRequestMethodKey:  "Method",
				},
			},
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assert.Equal(t, "basic error", evt.Message)
				assert.Equal(t, "svc#Method", evt.Transaction)
			},
		},
		{
			name: "simple error with http opts",
			err:  errors.New("basic error"),
			ropts: reporter.ReportOptions{
				Tags: map[string]interface{}{
					reporter.HTTPRequestMethodKey: "GET",
					reporter.HTTPRequestPathKey:   "/foo",
					reporter.HTTPRequestHostKey:   "example.com",
				},
			},
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assert.Equal(t, "basic error", evt.Message)
				assert.Equal(t, "GET /foo", evt.Transaction)
				assert.Equal(
					t,
					&sentry.Request{URL: "http://example.com/foo", Method: "GET"},
					evt.Request,
				)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewReporter(tt.opts...)

			assert.NoError(t, err)

			evt := r.buildEvent(tt.err, tt.ropts)
			tt.evtfn(t, evt)
		})
	}
}
