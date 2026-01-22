package subscriber

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/eventbus"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry/metric"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/event"
)

type metricSubscriber struct {
	telemetry telemetry.Telemetry
}

func NewMetricSubscriber(telemetry telemetry.Telemetry) eventbus.EventSubscriber {
	return &metricSubscriber{telemetry: telemetry}
}

func (s *metricSubscriber) SubscribedEvents() (syncHandlers eventbus.EventHandlersMap, asyncHandlers eventbus.EventHandlersMap) {
	return eventbus.EventHandlersMap{
		event.LeadCreatedEventName: s.TrackLeadCreated,
	}, nil
}

func (s *metricSubscriber) TrackLeadCreated(ctx context.Context, e domain.Event) {
	s.telemetry.RecordMetric(ctx, metric.LeadCounter{})
}
