// Package inhibit provides a reporter wrapper that can selectively suppress error reports.
//
// This package allows you to wrap any reporter.Reporter with error inhibition logic.
// Error inhibitors can inspect errors and decide whether they should be reported,
// which is useful for filtering out known errors, noise, or errors below certain thresholds.
package inhibit

import (
	"sync"

	"github.com/upfluence/errors/reporter"
)

// ErrorInhibitor determines whether an error should be inhibited from reporting.
type ErrorInhibitor interface {
	Inhibit(error) bool
}

// ErrorInhibitorFunc is a function adapter for the ErrorInhibitor interface.
type ErrorInhibitorFunc func(error) bool

func (fn ErrorInhibitorFunc) Inhibit(err error) bool { return fn(err) }

// Reporter wraps another reporter with error inhibition logic.
type Reporter struct {
	r reporter.Reporter

	mu  sync.RWMutex
	eis []ErrorInhibitor
}

// NewReporter creates a new inhibit reporter wrapping the given reporter
// with the specified error inhibitors.
func NewReporter(r reporter.Reporter, eis ...ErrorInhibitor) *Reporter {
	return &Reporter{r: r, eis: eis}
}

// AddErrorInhibitors adds additional error inhibitors to the reporter.
func (r *Reporter) AddErrorInhibitors(eis ...ErrorInhibitor) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.eis = append(r.eis, eis...)
}

// Close closes the underlying reporter.
func (r *Reporter) Close() error { return r.r.Close() }

// Report reports an error if it is not inhibited by any error inhibitors.
func (r *Reporter) Report(err error, opts reporter.ReportOptions) {
	r.mu.RLock()

	for _, ei := range r.eis {
		if ei.Inhibit(err) {
			r.mu.RUnlock()
			return
		}
	}

	r.mu.RUnlock()

	r.r.Report(err, opts)
}
