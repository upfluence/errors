package sentry

import "github.com/getsentry/sentry-go"

type tagIntegration struct {
	tags map[string]string
}

func (ti *tagIntegration) Name() string { return "custom_tags" }
func (ti *tagIntegration) SetupOnce(cl *sentry.Client) {
	cl.AddEventProcessor(ti.process)
}

func (ti *tagIntegration) process(evt *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	if evt.Tags == nil {
		evt.Tags = make(map[string]string, len(ti.tags))
	}

	for k, v := range ti.tags {
		evt.Tags[k] = v
	}

	return evt
}
