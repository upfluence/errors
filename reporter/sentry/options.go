package sentry

import (
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/upfluence/errors/reporter"
)

var defaultOptions = Options{
	SentryOptions: sentry.ClientOptions{Dsn: os.Getenv("SENTRY_DSN")},
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
	SentryOptions sentry.ClientOptions
	Timeout       time.Duration

	TagWhitelist map[string]struct{}
	TagBlacklist []func(string) bool
}

func (o Options) client() (*sentry.Client, error) {
	return sentry.NewClient(o.SentryOptions)
}

type Option func(*Options)
