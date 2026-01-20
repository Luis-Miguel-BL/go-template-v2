package event

import "github.com/Luis-Miguel-BL/go-lm-template/internal/domain"

const LeadCreatedEventName domain.EventName = "LeadCreated"

type LeadCreated struct {
	*domain.EventBase
	LeadUUID       string
	Name           string
	Email          string
	DocumentNumber string
	PhoneNumber    string
}

func (e LeadCreated) EventName() domain.EventName {
	return LeadCreatedEventName
}
