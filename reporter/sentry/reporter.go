package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/upfluence/errors/reporter"
)

type Reporter struct {
	cl *sentry.Client

	timeout time.Duration
}

func (r *Reporter) Report(err error, opts reporter.ReportOptions) {
}

func (r *Reporter) Close() error {
	r.cl.Flush(r.timeout)
	return nil
}
