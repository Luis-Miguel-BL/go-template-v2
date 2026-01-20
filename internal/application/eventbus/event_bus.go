package eventbus

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

type EventHandler func(context.Context, domain.Event)

type EventHandlersMap map[domain.EventName]EventHandler

type EventBus interface {
	Publish(ctx context.Context, event domain.Event)
}

type EventSubscriber interface {
	SubscribedEvents() (syncHandlers EventHandlersMap, asyncHandlers EventHandlersMap)
}
