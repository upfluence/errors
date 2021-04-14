package inhibit

import (
	"sync"

	"github.com/upfluence/errors/reporter"
)

type ErrorInhibitor interface {
	Inhibit(error) bool
}

type ErrorInhibitorFunc func(error) bool

func (fn ErrorInhibitorFunc) Inhibit(err error) bool { return fn(err) }

type Reporter struct {
	r reporter.Reporter

	mu  sync.RWMutex
	eis []ErrorInhibitor
}

func NewReporter(r reporter.Reporter, eis ...ErrorInhibitor) *Reporter {
	return &Reporter{r: r, eis: eis}
}

func (r *Reporter) AddErrorInhibitors(eis ...ErrorInhibitor) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.eis = append(r.eis, eis...)
}

func (r *Reporter) Close() error { return r.r.Close() }

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
