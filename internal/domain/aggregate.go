package domain

import "sync"

type Aggregate interface {
	AppendEvent(event Event)
	GetAndClearUncommitedEvents() []Event
}

type AggregateBase struct {
	uncommittedEvents []Event
	mu                *sync.Mutex
}

func NewAggregateBase() *AggregateBase {
	return &AggregateBase{
		mu: &sync.Mutex{},
	}
}

func (a *AggregateBase) AppendEvent(event Event) {
	a.uncommittedEvents = append(a.uncommittedEvents, event)
}

func (a *AggregateBase) GetAndClearUncommitedEvents() []Event {
	a.mu.Lock()
	defer a.mu.Unlock()
	events := a.uncommittedEvents
	a.uncommittedEvents = nil
	return events
}
