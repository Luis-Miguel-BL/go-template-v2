package model

import (
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/common/vo"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/event"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/util"
)

type LeadUUID string

type Lead struct {
	*domain.AggregateBase
	LeadUUID       LeadUUID
	Name           vo.PersonName
	Email          vo.EmailAddress
	Phone          vo.PhoneNumber
	DocumentNumber vo.DocumentNumber
	MotherName     vo.PersonName
	BirthDate      time.Time
}

type NewLeadParams struct {
	Name           vo.PersonName
	Email          vo.EmailAddress
	Phone          vo.PhoneNumber
	DocumentNumber vo.DocumentNumber
}

func NewLead(params NewLeadParams) *Lead {
	leadUUID := util.NewUUID()
	lead := &Lead{
		LeadUUID:       LeadUUID(leadUUID),
		AggregateBase:  domain.NewAggregateBase(),
		Name:           params.Name,
		Email:          params.Email,
		Phone:          params.Phone,
		DocumentNumber: params.DocumentNumber,
	}

	lead.AppendEvent(event.LeadCreated{
		EventBase:      domain.NewEventBase(),
		LeadUUID:       leadUUID,
		Name:           lead.Name.String(),
		Email:          lead.Email.String(),
		DocumentNumber: lead.DocumentNumber.String(),
		PhoneNumber:    lead.Phone.String(),
	})

	return lead
}

func (l *Lead) UpdateBirthDate(birthDate time.Time) {
	l.BirthDate = birthDate

	l.AppendEvent(event.LeadBirthDateUpdated{
		EventBase: domain.NewEventBase(),
		LeadUUID:  string(l.LeadUUID),
		BirthDate: birthDate,
	})
}

func (l *Lead) UpdateMotherName(motherName vo.PersonName) {
	l.MotherName = motherName

	l.AppendEvent(event.LeadMotherNameUpdated{
		EventBase:  domain.NewEventBase(),
		LeadUUID:   string(l.LeadUUID),
		MotherName: motherName.String(),
	})
}
