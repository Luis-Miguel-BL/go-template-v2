package event

import (
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

const LeadBirthDateUpdatedEventName domain.EventName = "LeadBirthDateUpdated"

type LeadBirthDateUpdated struct {
	*domain.EventBase
	LeadUUID  string
	BirthDate time.Time
}

func (e LeadBirthDateUpdated) EventName() domain.EventName {
	return LeadBirthDateUpdatedEventName
}
