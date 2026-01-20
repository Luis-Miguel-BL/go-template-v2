package subscriber

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/eventbus"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability/metric"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/event"
)

type metricSubscriber struct {
}

func NewMetricSubscriber() eventbus.EventSubscriber {
	return &metricSubscriber{}
}

func (s *metricSubscriber) SubscribedEvents() (syncHandlers eventbus.EventHandlersMap, asyncHandlers eventbus.EventHandlersMap) {
	return eventbus.EventHandlersMap{
		event.LeadCreatedEventName: s.TrackLeadCreated,
	}, nil
}

func (s *metricSubscriber) TrackLeadCreated(ctx context.Context, e domain.Event) {
	observability.GetObservability().RecordMetric(ctx, metric.LeadCounter{})
}
