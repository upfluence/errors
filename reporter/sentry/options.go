package sentry

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	uerrors "github.com/upfluence/errors"
	"github.com/upfluence/errors/base"
	"github.com/upfluence/errors/reporter"
)

type ErrorLevelFunc func(error) sentry.Level

var (
	defaultErrorLevelFuncs = []ErrorLevelFunc{
		ErrorIsFunc(context.DeadlineExceeded, sentry.LevelWarning),
		ErrorIsFunc(context.Canceled, sentry.LevelWarning),
		ErrorIsFunc(io.EOF, sentry.LevelWarning),
		ErrorCauseTextContainsFunc("net/http: TLS handshake timeout", sentry.LevelWarning),
		ErrorCauseTextContainsFunc("operation was canceled", sentry.LevelWarning),
		ErrorCauseTextContainsFunc("EOF", sentry.LevelWarning),
	}
	defaultOptions = Options{
		Tags: make(map[string]string),
		SentryOptions: sentry.ClientOptions{
			Dsn:            os.Getenv("SENTRY_DSN"),
			Environment:    os.Getenv("ENV"),
			SendDefaultPII: true,
		},
		TagWhitelist: toStringMap(
			[]string{reporter.RemoteIP, reporter.RemotePort, reporter.DomainKey},
		),
		Timeout: time.Minute,
		TagBlacklist: []func(string) bool{
			stringEqual(reporter.TransactionKey),
			stringEqual(reporter.UserEmailKey),
			stringEqual(reporter.UserIDKey),
			stringEqual(reporter.HTTPRequestProtoKey),
			stringEqual(reporter.HTTPRequestPathKey),
			stringEqual(reporter.HTTPRequestHostKey),
			stringEqual(reporter.HTTPRequestMethodKey),
			stringEqual(reporter.HTTPRequestBodyKey),
			stringPrefix(reporter.HTTPRequestHeaderKeyPrefix),
			stringPrefix(reporter.HTTPRequestQueryValuesKeyPrefix),
		},
		ErrorLevelFuncs: defaultErrorLevelFuncs,
	}
)

func stringEqual(x string) func(string) bool {
	return func(y string) bool { return x == y }
}

func stringPrefix(x string) func(string) bool {
	return func(y string) bool { return strings.HasPrefix(y, x) }
}

func toStringMap(vs []string) map[string]struct{} {
	res := make(map[string]struct{}, len(vs))

	for _, v := range vs {
		res[v] = struct{}{}
	}

	return res
}

// ErrorIsFunc creates an ErrorLevelFunc of the passed level that checks if
// reported errors are the same as the given sentinel error
func ErrorIsFunc(sentinel error, level sentry.Level) ErrorLevelFunc {
	return func(err error) sentry.Level {
		if uerrors.Is(err, sentinel) {
			return level
		}

		return ""
	}
}

// ErrorIsOfTypeFunc creates an ErrorLevelFunc of the passed level that checks
// if reported errors are of the passed generic type
func ErrorIsOfTypeFunc[T error](level sentry.Level) ErrorLevelFunc {
	return func(err error) sentry.Level {
		if uerrors.IsOfType[T](err) {
			return level
		}

		return ""
	}
}

// ErrorCauseTextContainsFunc creates an ErrorLevelFunc of the passed level that checks
// if reported errors' cause's Error() text contains the passed string
func ErrorCauseTextContainsFunc(errorText string, level sentry.Level) ErrorLevelFunc {
	return func(err error) sentry.Level {
		rootCause := base.UnwrapAll(err).Error()

		if strings.Contains(rootCause, errorText) {
			return level
		}

		return ""
	}
}

type Options struct {
	Tags map[string]string

	SentryOptions sentry.ClientOptions
	Timeout       time.Duration

	TagWhitelist map[string]struct{}
	TagBlacklist []func(string) bool

	ErrorLevelFuncs []ErrorLevelFunc
}

func (o Options) client() (*sentry.Client, error) {
	if len(o.Tags) > 0 {
		o.SentryOptions.Integrations = func(is []sentry.Integration) []sentry.Integration {
			return append(is, &tagIntegration{tags: o.Tags})
		}
	}

	return sentry.NewClient(o.SentryOptions)
}

type Option func(*Options)

// WithTags allows for custom tags to be given to the Reporter
func WithTags(tags map[string]string) Option {
	return func(opts *Options) {
		for k, v := range tags {
			opts.Tags[k] = v
		}
	}
}

// AppendErrorLevelFuncs adds the passed funcs to the ErrorLevelFuncs of the Reporter
func AppendErrorLevelFuncs(funcs ...ErrorLevelFunc) Option {
	return func(opts *Options) {
		opts.ErrorLevelFuncs = append(opts.ErrorLevelFuncs, funcs...)
	}
}

// ReplaceErrorLevelFuncs replaces the ErrorLevelFuncs of the Reporter with the passed ones
func ReplaceErrorLevelFuncs(funcs []ErrorLevelFunc) Option {
	return func(opts *Options) {
		opts.ErrorLevelFuncs = funcs
	}
}
