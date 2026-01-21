package subscriber

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/eventbus"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	obs_event "github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability/event"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/event"
)

type monitorSubscriber struct {
	obs observability.Observability
}

func NewMonitorSubscriber(obs observability.Observability) eventbus.EventSubscriber {
	return &monitorSubscriber{obs: obs}
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

	s.obs.AddEvent(ctx, obs_event.LeadCreated{
		LeadID:    event.LeadUUID,
		CreatedAt: event.OccurredAt(),
	})
}
