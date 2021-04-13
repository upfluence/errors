package sentry

import (
	"fmt"
	"go/build"
	"net/url"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/upfluence/errors/base"
	"github.com/upfluence/errors/domain"
	"github.com/upfluence/errors/reporter"
	"github.com/upfluence/errors/stacktrace"
	"github.com/upfluence/errors/tags"
)

type Reporter struct {
	cl *sentry.Client

	tagWhitelist map[string]struct{}
	tagBlacklist []func(string) bool

	timeout time.Duration
}

func NewReporter(os ...Option) (*Reporter, error) {
	var opts = defaultOptions

	for _, o := range os {
		o(&opts)
	}

	cl, err := opts.client()

	if err != nil {
		return nil, err
	}

	return &Reporter{
		cl:           cl,
		tagWhitelist: opts.TagWhitelist,
		tagBlacklist: opts.TagBlacklist,
		timeout:      opts.Timeout,
	}, nil
}

func (r *Reporter) Report(err error, opts reporter.ReportOptions) {
	evt := r.buildEvent(err, opts)

	if evt == nil {
		return
	}

	r.cl.CaptureEvent(evt, nil, nil)
}

func (r *Reporter) Close() error {
	r.cl.Flush(r.timeout)
	return nil
}

func (r *Reporter) appendTag(k string, v interface{}, evt *sentry.Event) {
	for _, fn := range r.tagBlacklist {
		if fn(k) {
			return
		}
	}

	if _, ok := r.tagWhitelist[k]; ok {
		evt.Tags[k] = stringifyTag(v)
		return
	}

	evt.Contexts[k] = v
}

func (r *Reporter) buildEvent(err error, opts reporter.ReportOptions) *sentry.Event {
	if err == nil {
		return nil
	}

	ts := tags.GetTags(err)

	for k, v := range opts.Tags {
		if _, ok := ts[k]; !ok {
			ts[k] = v
		}
	}

	evt := sentry.NewEvent()

	evt.Timestamp = time.Now()
	evt.Message = err.Error()
	evt.Transaction = transactionName(ts)
	evt.User = buildUser(ts)
	evt.Request = buildRequest(ts)

	cause := base.UnwrapAll(err)

	evt.Exception = []sentry.Exception{
		{
			Type:       fmt.Sprintf("%T", cause),
			Value:      cause.Error(),
			Module:     string(domain.GetDomain(err)),
			Stacktrace: extractStacktrace(err, opts.Depth+2),
		},
	}

	for k, v := range ts {
		r.appendTag(k, v, evt)
	}

	return evt
}

func extractStacktrace(err error, n int) *sentry.Stacktrace {
	var s sentry.Stacktrace

	appendFrame := func(f stacktrace.Frame) {
		fn, file, line := f.Location()

		pkg, fn := splitQualifiedFunctionName(fn)

		s.Frames = append(
			s.Frames,
			sentry.Frame{
				AbsPath:  file,
				Function: fn,
				Lineno:   line,
				Module:   pkg,
				InApp:    isInApp(file, pkg),
			},
		)
	}

	appendFrame(stacktrace.Caller(n + 1))

	for _, f := range stacktrace.GetFrames(err) {
		appendFrame(f)
	}

	return &s
}

func splitQualifiedFunctionName(name string) (string, string) {
	pkg := stacktrace.PackageName(name)
	return pkg, strings.TrimPrefix(name, pkg+".")
}

func isInApp(absPath, module string) bool {
	if strings.HasPrefix(absPath, build.Default.GOROOT) ||
		strings.Contains(module, "vendor") ||
		strings.Contains(module, "third_party") {
		return false
	}

	return true
}

func stringifyTag(v interface{}) string {
	if v == nil {
		return ""
	}

	switch vv := v.(type) {
	case string:
		return vv
	case []byte:
		return string(vv)
	case fmt.Stringer:
		return vv.String()
	}

	return fmt.Sprint(v)
}

func transactionName(tags map[string]interface{}) string {
	if v, ok := tags[reporter.TransactionKey]; ok {
		return stringifyTag(v)
	}

	if v, ok := tags[reporter.ThriftRequestMethodKey]; ok {
		return fmt.Sprintf(
			"%s#%s",
			stringifyTag(tags[reporter.ThriftRequestServiceKey]),
			stringifyTag(v),
		)
	}

	if v, ok := tags[reporter.HTTPRequestPathKey]; ok {
		return fmt.Sprintf(
			"%s %s",
			stringifyTag(tags[reporter.HTTPRequestMethodKey]),
			stringifyTag(v),
		)
	}

	return ""
}

func buildUser(tags map[string]interface{}) sentry.User {
	return sentry.User{
		Email: stringifyTag(tags[reporter.UserEmailKey]),
		ID:    stringifyTag(tags[reporter.UserIDKey]),
	}
}

func buildRequest(tags map[string]interface{}) *sentry.Request {
	var (
		req    sentry.Request
		tapped bool

		u   = url.URL{Scheme: "http", Host: "localhost"}
		qvs = url.Values{}
	)

	for k, v := range tags {
		switch {
		case k == reporter.HTTPRequestProtoKey:
			u.Scheme = stringifyTag(v)
		case k == reporter.HTTPRequestPathKey:
			u.Path = stringifyTag(v)
		case k == reporter.HTTPRequestHostKey:
			u.Host = stringifyTag(v)
		case k == reporter.HTTPRequestMethodKey:
			req.Method = stringifyTag(v)
		case strings.HasPrefix(k, reporter.HTTPRequestHeaderKeyPrefix):
			if req.Headers == nil {
				req.Headers = make(map[string]string)
			}

			k := strings.TrimPrefix(k, reporter.HTTPRequestHeaderKeyPrefix)
			req.Headers[k] = stringifyTag(v)
		case strings.HasPrefix(k, reporter.HTTPRequestQueryValuesKeyPrefix):
			qvs.Add(
				strings.TrimPrefix(k, reporter.HTTPRequestQueryValuesKeyPrefix),
				stringifyTag(v),
			)
		default:
			continue
		}

		tapped = true
	}

	if tapped {
		req.URL = u.String()
		req.QueryString = qvs.Encode()

		return &req
	}

	return nil
}
