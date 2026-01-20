package event

import (
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

const LeadMotherNameUpdatedEventName domain.EventName = "LeadMotherNameUpdated"

type LeadMotherNameUpdated struct {
	*domain.EventBase
	LeadUUID   string
	MotherName string
}

func (e LeadMotherNameUpdated) EventName() domain.EventName {
	return LeadMotherNameUpdatedEventName
}
