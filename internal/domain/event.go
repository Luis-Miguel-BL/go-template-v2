package domain

import (
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/util"
)

type EventName string
type EventID string

type Event interface {
	EventName() EventName
	EventID() EventID
	OccurredAt() time.Time
}

type EventBase struct {
	eventID    EventID
	occurredAt time.Time
}

func NewEventBase() *EventBase {
	return &EventBase{
		eventID:    EventID(util.NewUUID()),
		occurredAt: time.Now(),
	}
}

func (e *EventBase) EventID() EventID {
	return e.eventID
}
func (e *EventBase) OccurredAt() time.Time {
	return e.occurredAt
}
