package domain

import (
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/util"
)

type EventName string
type EventID string
type EventMetadata struct {
}

type Event interface {
	EventName() EventName
	EventID() EventID
	OccurredAt() time.Time
	Metadata() EventMetadata
}

type EventBase struct {
	eventID    EventID
	occurredAt time.Time
	metadata   EventMetadata
}

func NewEventBase() *EventBase {
	return &EventBase{
		eventID:    EventID(util.NewUUID()),
		occurredAt: time.Now(),
		metadata:   EventMetadata{},
	}
}

func (e *EventBase) WithMetadata(metadata EventMetadata) *EventBase {
	e.metadata = metadata
	return e
}

func (e *EventBase) EventID() EventID {
	return e.eventID
}
func (e *EventBase) OccurredAt() time.Time {
	return e.occurredAt
}
func (e *EventBase) Metadata() EventMetadata {
	return e.metadata
}
