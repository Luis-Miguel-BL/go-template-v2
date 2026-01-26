package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/auth"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/auth/mocks"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/dto"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	jwtMock *mocks.JWTHelper[TokenClaims]
	service *authService
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.jwtMock = new(mocks.JWTHelper[TokenClaims])

	cfg := &config.Config{
		Server: config.ServerConfig{
			AppKey: "my-secret-key",
		},
	}

	s.service = NewAuthService(s.jwtMock, cfg)
}

func (s *AuthServiceTestSuite) TearDownTest() {
	s.jwtMock.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestGenerateToken_Success() {
	s.jwtMock.
		On("GenerateToken", mock.Anything, time.Hour*24, mock.Anything).
		Return("token123", int64(3600), nil)

	token, err := s.service.GenerateToken(context.Background())

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "token123", token.Token)
	assert.Equal(s.T(), int64(3600), token.ExpiresIn)
}

func (s *AuthServiceTestSuite) TestGenerateToken_Error() {
	s.jwtMock.
		On("GenerateToken", mock.Anything, time.Hour*24, mock.Anything).
		Return("", int64(0), errors.New("jwt error"))

	token, err := s.service.GenerateToken(context.Background())

	assert.Error(s.T(), err)
	assert.Equal(s.T(), dto.AccessToken{}, token)
}

func (s *AuthServiceTestSuite) TestValidateToken_Success() {
	claims := &TokenClaims{
		LeadID:    ptr("lead-1"),
		SessionID: ptr("session-1"),
	}

	s.jwtMock.
		On("ValidateToken", "abc").
		Return(claims, nil)

	result, err := s.service.ValidateToken(context.Background(), "abc")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), claims, result)
}

func (s *AuthServiceTestSuite) TestRefreshToken_PreserveClaims() {
	ctx := auth.WithContext(context.Background(), &TokenClaims{})

	inputClaims := &TokenClaims{
		LeadID:    ptr("lead-x"),
		SessionID: ptr("session-y"),
	}

	s.jwtMock.
		On("GenerateToken", mock.Anything, time.Hour*24, TokenClaims{
			LeadID:    inputClaims.LeadID,
			SessionID: inputClaims.SessionID,
		}).
		Return("newtoken", int64(7200), nil)

	token, err := s.service.RefreshToken(ctx, inputClaims)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "newtoken", token.Token)
	assert.Equal(s.T(), int64(7200), token.ExpiresIn)
}

func (s *AuthServiceTestSuite) TestValidateAppKey() {
	assert.True(s.T(), s.service.ValidateAppKey(context.Background(), "my-secret-key"))
	assert.False(s.T(), s.service.ValidateAppKey(context.Background(), "wrong-key"))
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func ptr(s string) *string {
	return &s
}
