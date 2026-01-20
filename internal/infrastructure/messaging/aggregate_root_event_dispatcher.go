package messaging

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/eventbus"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

type AggregateRootEventDispatcher struct {
	eventbus eventbus.EventBus
}

func NewAggregateRootEventDispatcher(eventbus eventbus.EventBus) *AggregateRootEventDispatcher {
	return &AggregateRootEventDispatcher{
		eventbus: eventbus,
	}
}

func (d *AggregateRootEventDispatcher) PublishUncommitedEvents(ctx context.Context, aggregate domain.Aggregate) {
	uncommittedEvents := aggregate.GetAndClearUncommitedEvents()

	for _, event := range uncommittedEvents {
		d.eventbus.Publish(ctx, event)
	}
}
