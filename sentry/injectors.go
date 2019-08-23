package sentry

import (
	"github.com/getsentry/sentry-go"
)

type Injector interface {
	Inject(event *sentry.Event) *sentry.Event
}

type TagInjector struct {
	Tags func() map[string]string
}

func (injector TagInjector) Inject(event *sentry.Event) *sentry.Event {
	if event.Tags == nil {
		event.Tags = make(map[string]string)
	}

	tags := injector.Tags()
	for k, v := range tags {
		event.Tags[k] = v
	}

	return event
}
