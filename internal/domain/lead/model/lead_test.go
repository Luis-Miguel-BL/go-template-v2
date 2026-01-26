package model

import (
	"testing"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/common/vo"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LeadTestSuite struct {
	suite.Suite
	lead *Lead
}

func (s *LeadTestSuite) SetupTest() {
	s.lead = NewLead(validParams())
}

func validParams() NewLeadParams {
	name, _ := vo.NewPersonName("John Doe")
	email, _ := vo.NewEmailAddress("john@mail.com")
	phone, _ := vo.NewPhoneNumber("31999999999")
	documentNumber, _ := vo.NewDocumentNumber("12345678909")
	return NewLeadParams{
		Name:           name,
		Email:          email,
		Phone:          phone,
		DocumentNumber: documentNumber,
	}
}

func (s *LeadTestSuite) TestNewLead() {

	s.Run("Should create lead with correct fields", func() {
		lead := NewLead(validParams())

		assert.NotNil(s.T(), lead)
		assert.NotEmpty(s.T(), lead.LeadUUID)
		assert.Equal(s.T(), validParams().Name, lead.Name)
		assert.Equal(s.T(), validParams().Email, lead.Email)
		assert.Equal(s.T(), validParams().Phone, lead.Phone)
		assert.Equal(s.T(), validParams().DocumentNumber, lead.DocumentNumber)
	})

	s.Run("Should raise LeadCreated event", func() {

		lead := NewLead(validParams())

		events := lead.GetAndClearUncommitedEvents()
		assert.Len(s.T(), events, 1)

		evt, ok := events[0].(event.LeadCreated)
		assert.True(s.T(), ok)

		assert.Equal(s.T(), string(lead.LeadUUID), evt.LeadUUID)
		assert.Equal(s.T(), lead.Name.String(), evt.Name)
		assert.Equal(s.T(), lead.Email.String(), evt.Email)
		assert.Equal(s.T(), lead.Phone.String(), evt.PhoneNumber)
		assert.Equal(s.T(), lead.DocumentNumber.String(), evt.DocumentNumber)
	})
}

func (s *LeadTestSuite) TestUpdateBirthDate() {

	s.Run("Should update birth date", func() {
		date := time.Date(1995, 3, 10, 0, 0, 0, 0, time.UTC)

		s.lead.UpdateBirthDate(date)

		assert.Equal(s.T(), date, s.lead.BirthDate)
	})

	s.Run("Should raise LeadBirthDateUpdated event", func() {

		s.lead.GetAndClearUncommitedEvents()

		date := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

		s.lead.UpdateBirthDate(date)

		events := s.lead.GetAndClearUncommitedEvents()
		assert.Len(s.T(), events, 1)

		evt, ok := events[0].(event.LeadBirthDateUpdated)
		assert.True(s.T(), ok)

		assert.Equal(s.T(), string(s.lead.LeadUUID), evt.LeadUUID)
		assert.Equal(s.T(), date, evt.BirthDate)
	})
}

func (s *LeadTestSuite) TestUpdateMotherName() {

	s.Run("Should update mother name", func() {

		mother, err := vo.NewPersonName("Mary Doe")
		assert.NoError(s.T(), err)

		s.lead.UpdateMotherName(mother)

		assert.Equal(s.T(), mother, s.lead.MotherName)
	})

	s.Run("Should raise LeadMotherNameUpdated event", func() {

		s.lead.GetAndClearUncommitedEvents()

		mother, err := vo.NewPersonName("Jane Doe")
		assert.NoError(s.T(), err)

		s.lead.UpdateMotherName(mother)

		events := s.lead.GetAndClearUncommitedEvents()
		assert.Len(s.T(), events, 1)

		evt, ok := events[0].(event.LeadMotherNameUpdated)
		assert.True(s.T(), ok)

		assert.Equal(s.T(), string(s.lead.LeadUUID), evt.LeadUUID)
		assert.Equal(s.T(), mother.String(), evt.MotherName)
	})
}

func TestLeadTestSuite(t *testing.T) {
	suite.Run(t, new(LeadTestSuite))
}
