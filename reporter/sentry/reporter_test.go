package sentry

import (
	"io"
	"regexp"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"

	"github.com/upfluence/errors"
	"github.com/upfluence/errors/recovery"
	"github.com/upfluence/errors/reporter"
)

func TestBuildEvent(t *testing.T) {
	for _, tt := range []struct {
		name string

		opts      []Option
		modifiers []func(*Reporter)
		err       error
		ropts     reporter.ReportOptions

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
				assert.Equal(t, map[string]map[string]interface{}{}, evt.Contexts)
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
		{

			name: "simple error with http opts",
			err:  recovery.WrapRecoverResult(1),
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assertFuncNames(
					t,
					evt,
					[]string{`TestBuildEvent.func\d`, "tRunner", "goexit"},
				)
			},
		},
		{
			name: "simple error with http opts",
			err: errors.WithTags(
				errors.New("basic error"),
				map[string]interface{}{"foo": "bar", "biz": "buz"},
			),
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assert.Equal(t, evt.Extra, map[string]interface{}{"foo": "bar", "biz": "buz"})
				assert.Equal(t, evt.Tags, map[string]string{"domain": "github.com/upfluence/errors/reporter/sentry"})
			},
		},
		{
			name: "simple error with extra",
			err: errors.WithTags(
				errors.New("basic error"),
				map[string]interface{}{"foo": "bar", "biz": "buz"},
			),
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assert.Equal(t, evt.Extra, map[string]interface{}{"foo": "bar", "biz": "buz"})
				assert.Equal(t, evt.Tags, map[string]string{"domain": "github.com/upfluence/errors/reporter/sentry"})
			},
		},
		{
			name: "simple error with extra & whitelisted tag",
			err: errors.WithTags(
				errors.New("basic error"),
				map[string]interface{}{"foo": "bar", "biz": "buz"},
			),
			modifiers: []func(*Reporter){
				func(r *Reporter) {
					r.WhitelistTag(func(s string) bool { return s == "foo" })
				},
			},
			evtfn: func(t *testing.T, evt *sentry.Event) {
				assert.Equal(t, evt.Extra, map[string]interface{}{"biz": "buz"})
				assert.Equal(t, evt.Tags, map[string]string{"foo": "bar", "domain": "github.com/upfluence/errors/reporter/sentry"})
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewReporter(tt.opts...)

			assert.NoError(t, err)

			for _, fn := range tt.modifiers {
				fn(r)
			}

			evt := r.buildEvent(tt.err, tt.ropts)
			tt.evtfn(t, evt)
		})
	}
}

func TestSegfault(t *testing.T) {
	r, err := NewReporter()

	assert.NoError(t, err)

	defer func() {
		if err := recover(); err != nil {
			evt := r.buildEvent(
				recovery.WrapRecoverResult(err),
				reporter.ReportOptions{},
			)

			assertFuncNames(
				t,
				evt,
				[]string{`TestSegfault.func\d`, "gopanic", "panicmem", "sigpanic", "TestSegfault", "tRunner", "goexit"},
			)
		} else {
			t.Error("no error recovered")
		}
	}()

	var f io.Reader

	f.Read(nil)
}

func assertFuncNames(t testing.TB, evt *sentry.Event, want []string) {
	exc := evt.Exception
	assert.Len(t, exc, 1)

	frames := exc[0].Stacktrace.Frames

	assert.Len(t, frames, len(want))

	for i, f := range frames {
		assert.True(
			t,
			regexp.MustCompile(want[i]).MatchString(f.Function),
			"want regexp = %s, function = %s",
			want[i],
			f.Function,
		)
	}
}
