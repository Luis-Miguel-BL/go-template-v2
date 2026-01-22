package subscriber

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/eventbus"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	telemetry_event "github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry/event"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/event"
)

type monitorSubscriber struct {
	telemetry telemetry.Telemetry
}

func NewMonitorSubscriber(telemetry telemetry.Telemetry) eventbus.EventSubscriber {
	return &monitorSubscriber{telemetry: telemetry}
}

func (s *monitorSubscriber) SubscribedEvents() (syncHandlers eventbus.EventHandlersMap, asyncHandlers eventbus.EventHandlersMap) {
	return eventbus.EventHandlersMap{
		event.LeadCreatedEventName: s.TrackLeadCreated,
	}, nil
}

func (s *monitorSubscriber) TrackLeadCreated(ctx context.Context, e domain.Event) {
	log := logger.FromContext(ctx)
	event := e.(event.LeadCreated)
	log.Debug("Lead created")

	s.telemetry.AddEvent(ctx, telemetry_event.LeadCreated{
		LeadID:    event.LeadUUID,
		CreatedAt: event.OccurredAt(),
	})
}
