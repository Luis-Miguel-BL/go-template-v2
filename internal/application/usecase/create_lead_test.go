package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/dto"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/service/mocks"
	leadmocks "github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateLeadTestSuite struct {
	suite.Suite
	authService *mocks.AuthService
	leadRepo    *leadmocks.LeadRepository
	usecase     *CreateLead
}

func (s *CreateLeadTestSuite) SetupTest() {
	s.authService = new(mocks.AuthService)
	s.leadRepo = new(leadmocks.LeadRepository)
	s.usecase = NewCreateLead(s.authService, s.leadRepo)
}

func (s *CreateLeadTestSuite) TearDownTest() {
	s.authService.AssertExpectations(s.T())
	s.leadRepo.AssertExpectations(s.T())
}

func validInput() dto.CreateLeadInput {
	return dto.CreateLeadInput{
		Name:           "John Doe",
		Email:          "john@mail.com",
		Phone:          "31999999999",
		DocumentNumber: "12345678909",
	}
}

func (s *CreateLeadTestSuite) TestExecute() {

	s.Run("Should create lead and return token when input is valid", func() {
		input := validInput()

		s.leadRepo.
			On("Save", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		s.authService.
			On("RefreshToken", mock.Anything, mock.Anything).
			Return(dto.AccessToken{
				Token:     "token123",
				ExpiresIn: 3600,
			}, nil).
			Once()

		output, err := s.usecase.Execute(context.Background(), input)

		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), output.LeadID)
		assert.Equal(s.T(), "token123", output.AccessToken.Token)
	})

	s.Run("Should return error when name is invalid", func() {
		input := validInput()
		input.Name = ""

		output, err := s.usecase.Execute(context.Background(), input)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), dto.CreateLeadOutput{}, output)
	})

	s.Run("Should return error when repository save fails", func() {
		input := validInput()

		s.leadRepo.
			On("Save", mock.Anything, mock.Anything).
			Return(errors.New("db error")).
			Once()

		output, err := s.usecase.Execute(context.Background(), input)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), dto.CreateLeadOutput{}, output)
	})

	s.Run("Should return error when refresh token fails", func() {
		input := validInput()

		s.leadRepo.
			On("Save", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		s.authService.
			On("RefreshToken", mock.Anything, mock.Anything).
			Return(dto.AccessToken{}, errors.New("jwt error")).
			Once()

		output, err := s.usecase.Execute(context.Background(), input)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), dto.CreateLeadOutput{}, output)
	})
}

func TestCreateLeadTestSuite(t *testing.T) {
	suite.Run(t, new(CreateLeadTestSuite))
}
