package messaging

import (
	"context"
	"reflect"
	"runtime"
	"sync"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/eventbus"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

type domainEventBus struct {
	wg            *sync.WaitGroup
	syncHandlers  map[domain.EventName][]eventbus.EventHandler
	asyncHandlers map[domain.EventName][]eventbus.EventHandler
	handlerNames  map[uintptr]string
}

func NewDomainEventBus(wg *sync.WaitGroup, eventSubscribers ...eventbus.EventSubscriber) eventbus.EventBus {
	bus := &domainEventBus{
		wg:            wg,
		syncHandlers:  make(map[domain.EventName][]eventbus.EventHandler),
		asyncHandlers: make(map[domain.EventName][]eventbus.EventHandler),
		handlerNames:  make(map[uintptr]string),
	}
	bus.subscribeAll(eventSubscribers...)
	return bus
}

func (b *domainEventBus) Publish(ctx context.Context, event domain.Event) {
	asyncHandlersCopy := b.asyncHandlers
	syncHandlersCopy := b.syncHandlers

	topic := event.EventName()

	if handlers, ok := asyncHandlersCopy[topic]; ok {
		for _, handler := range handlers {
			b.wg.Add(1)
			go b.doPublish(handler, ctx, event)
		}
	}

	if handlers, ok := syncHandlersCopy[topic]; ok {
		for _, handler := range handlers {
			b.wg.Add(1)
			b.doPublish(handler, ctx, event)
		}
	}
}

func getHandlerPointer(handler eventbus.EventHandler) uintptr {
	return reflect.ValueOf(handler).Pointer()
}

func getHandlerName(handler eventbus.EventHandler) string {
	return runtime.FuncForPC(getHandlerPointer(handler)).Name()
}

func (b *domainEventBus) subscribeAll(eventSubscribers ...eventbus.EventSubscriber) {
	for _, subscriber := range eventSubscribers {
		syncHandlers, asyncHandlers := subscriber.SubscribedEvents()

		for topic, handler := range syncHandlers {
			b.syncHandlers[topic] = append(b.syncHandlers[topic], handler)
			b.handlerNames[getHandlerPointer(handler)] = getHandlerName(handler)
		}
		for topic, handler := range asyncHandlers {
			b.asyncHandlers[topic] = append(b.asyncHandlers[topic], handler)
			b.handlerNames[getHandlerPointer(handler)] = getHandlerName(handler)
		}
	}
}

func (b *domainEventBus) doPublish(handler eventbus.EventHandler, ctx context.Context, event domain.Event) {
	defer b.wg.Done()
	ctx, span := observability.GetObservability().StartSpan(ctx, "EventBus.Publish")
	defer span.End()

	handlerName := b.handlerNames[getHandlerPointer(handler)]

	span.SetAttributes(map[string]any{
		"event.name":   event.EventName(),
		"event.id":     event.EventID(),
		"handler.name": handlerName,
	})
	handler(ctx, event)
}
