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

type ErrorLevelMapper func(error) sentry.Level

var defaultErrorLevelMappers = []ErrorLevelMapper{
	ErrorIsLevel(context.DeadlineExceeded, sentry.LevelWarning),
	ErrorIsLevel(context.Canceled, sentry.LevelWarning),
	ErrorIsLevel(io.EOF, sentry.LevelWarning),
	ErrorCauseTextContainsLevel("net/http: TLS handshake timeout", sentry.LevelWarning),
	ErrorCauseTextContainsLevel("operation was canceled", sentry.LevelWarning),
	ErrorCauseTextContainsLevel("EOF", sentry.LevelWarning),
}

func defaultOptions() Options {
	return Options{
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
		ErrorLevelMappers: defaultErrorLevelMappers,
	}
}

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

// ErrorIsLevel creates an ErrorLevelMapper of the passed level that checks if
// reported errors are the same as the given sentinel error
func ErrorIsLevel(sentinel error, level sentry.Level) ErrorLevelMapper {
	return func(err error) sentry.Level {
		if uerrors.Is(err, sentinel) {
			return level
		}

		return ""
	}
}

// ErrorIsOfTypeLevel creates an ErrorLevelMapper of the passed level that checks
// if reported errors are of the passed generic type
func ErrorIsOfTypeLevel[T error](level sentry.Level) ErrorLevelMapper {
	return func(err error) sentry.Level {
		if uerrors.IsOfType[T](err) {
			return level
		}

		return ""
	}
}

// ErrorCauseTextContainsLevel creates an ErrorLevelMapper of the passed level that checks
// if reported errors' cause's Error() text contains the passed string
func ErrorCauseTextContainsLevel(errorText string, level sentry.Level) ErrorLevelMapper {
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

	ErrorLevelMappers []ErrorLevelMapper
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

// AppendErrorLevelMappers adds the passed funcs to the ErrorLevelMappers of the Reporter
func AppendErrorLevelMappers(funcs ...ErrorLevelMapper) Option {
	return func(opts *Options) {
		opts.ErrorLevelMappers = append(opts.ErrorLevelMappers, funcs...)
	}
}

// ReplaceErrorLevelMappers replaces the ErrorLevelMappers of the Reporter with the passed ones
func ReplaceErrorLevelMappers(funcs []ErrorLevelMapper) Option {
	return func(opts *Options) {
		opts.ErrorLevelMappers = funcs
	}
}
