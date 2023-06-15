package sentry

import (
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/upfluence/errors/reporter"
)

var defaultOptions = Options{
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

type Options struct {
	Tags map[string]string

	SentryOptions sentry.ClientOptions
	Timeout       time.Duration

	TagWhitelist map[string]struct{}
	TagBlacklist []func(string) bool
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

func WithTags(tags map[string]string) Option {
	return func(opts *Options) {
		for k, v := range tags {
			opts.Tags[k] = v
		}
	}
}
