// Package sentry provides a Sentry implementation of the reporter.Reporter interface.
//
// This package integrates with Sentry for error tracking and monitoring.
// It automatically extracts error metadata including stacktraces, tags, domains,
// and request information, formatting them for Sentry ingestion. The reporter
// supports tag whitelisting/blacklisting for controlling which metadata is sent.
package sentry

import (
	"errors"
	"fmt"
	"go/build"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/upfluence/errors/base"
	"github.com/upfluence/errors/domain"
	"github.com/upfluence/errors/reporter"
	"github.com/upfluence/errors/stacktrace"
	"github.com/upfluence/errors/tags"
)

// Reporter is a Sentry error reporter implementation.
type Reporter struct {
	cl *sentry.Client

	tagWhitelist    []func(string) bool
	tagBlacklist    []func(string) bool
	errorLevelFuncs []ErrorLevelFunc

	timeout time.Duration
}

// NewReporter creates a new Sentry reporter with the given options.
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
		cl: cl,
		tagWhitelist: []func(string) bool{
			func(k string) bool {
				_, ok := opts.TagWhitelist[k]
				return ok
			},
		},
		tagBlacklist:    opts.TagBlacklist,
		timeout:         opts.Timeout,
		errorLevelFuncs: opts.ErrorLevelFuncs,
	}, nil
}

// WhitelistTag adds tag whitelist functions that determine which tags
// should be included as Sentry tags (vs extra data).
func (r *Reporter) WhitelistTag(fns ...func(string) bool) {
	r.tagWhitelist = append(r.tagWhitelist, fns...)
}

// Report sends an error to Sentry with the given options.
func (r *Reporter) Report(err error, opts reporter.ReportOptions) {
	evt := r.buildEvent(err, opts)

	fmt.Println(evt)

	if evt == nil {
		return
	}

	fmt.Println(*r.cl.CaptureEvent(evt, nil, nil))
}

// Close flushes pending events to Sentry and releases resources.
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

	for _, fn := range r.tagWhitelist {
		if fn(k) {
			evt.Tags[k] = stringifyTag(v)
			return
		}
	}

	evt.Extra[k] = v
}

func buildErrorTypeChain(err error) []string {
	var res []string

	for err != nil {
		res = append(res, fmt.Sprintf("%T", err))
		err = base.UnwrapOnce(err)
	}

	if len(res) < 2 {
		return nil
	}

	return res
}

func (r *Reporter) buildEvent(err error, opts reporter.ReportOptions) *sentry.Event {
	if err == nil {
		return nil
	}

	errorTags := tags.GetTags(err)

	if errorTags == nil && len(opts.Tags) > 0 {
		errorTags = make(map[string]interface{}, len(opts.Tags))
	}

	for k, v := range opts.Tags {
		if _, ok := errorTags[k]; !ok {
			errorTags[k] = v
		}
	}

	evt := sentry.NewEvent()

	evt.Level = r.computeErrorLevel(err)
	evt.Timestamp = time.Now()
	evt.Message = err.Error()
	evt.Transaction = transactionName(errorTags)
	evt.User = buildUser(errorTags)
	evt.Request = buildRequest(errorTags)

	cause := base.UnwrapAll(err)

	evt.Exception = []sentry.Exception{
		{
			Type:       fmt.Sprintf("%T", cause),
			Value:      cause.Error(),
			Module:     string(domain.GetDomain(err)),
			Stacktrace: extractStacktrace(err, opts.Depth+1),
		},
	}

	for k, v := range errorTags {
		r.appendTag(k, v, evt)
	}

	if ts := buildErrorTypeChain(err); len(ts) > 0 {
		r.appendTag("error_types", ts, evt)
	}

	return evt
}

var validLevels = []sentry.Level{
	sentry.LevelDebug,
	sentry.LevelInfo,
	sentry.LevelWarning,
	sentry.LevelError,
	sentry.LevelFatal,
}

func (r *Reporter) computeErrorLevel(err error) sentry.Level {
	for _, errFunc := range r.errorLevelFuncs {
		if level := errFunc(err); level != "" {
			return level
		}
	}

	return sentry.LevelError
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

	var rerr runtime.Error

	if errors.As(err, &rerr) {
		for _, f := range stacktrace.Stacktrace(n+1, 10) {
			appendFrame(f)
		}
	} else {
		appendFrame(stacktrace.Caller(n + 2))

		for _, f := range stacktrace.GetFrames(err) {
			appendFrame(f)
		}
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

	return fmt.Sprintf("%v", v)
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
		case k == reporter.HTTPRequestBodyKey:
			req.Data = stringifyTag(v)
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
