package service

import (
	"context"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/auth"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/dto"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/util"
)

type AuthService struct {
	jwtHelper       auth.JWTHelper[TokenClaims]
	appKey          string
	tokenExpiration time.Duration
}

func NewAuthService(jwtHelper auth.JWTHelper[TokenClaims], cfg *config.Config) *AuthService {
	return &AuthService{
		jwtHelper:       jwtHelper,
		appKey:          cfg.Server.AppKey,
		tokenExpiration: time.Hour * 24,
	}
}

type TokenClaims struct {
	LeadID    *string `json:"lead_id"`
	SessionID *string `json:"session_id"`
}

func (s *AuthService) GenerateToken(ctx context.Context) (accessToken dto.AccessToken, err error) {
	sessionID := util.NewUUID()
	clains := TokenClaims{
		SessionID: &sessionID,
	}

	token, expiresIn, err := s.jwtHelper.GenerateToken(ctx, s.tokenExpiration, clains)
	if err != nil {
		return accessToken, err
	}

	return dto.AccessToken{
		Token:     token,
		ExpiresIn: expiresIn,
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (tokenClaims *TokenClaims, err error) {
	return s.jwtHelper.ValidateToken(token)
}

func (s *AuthService) RefreshToken(ctx context.Context, claims *TokenClaims) (newAccessToken dto.AccessToken, err error) {
	newClaims := auth.FromContext[TokenClaims](ctx)

	if claims != nil {
		if claims.LeadID != nil {
			newClaims.LeadID = claims.LeadID
		}
		if claims.SessionID != nil {
			newClaims.SessionID = claims.SessionID
		}
	}
	token, expiresIn, err := s.jwtHelper.GenerateToken(ctx, s.tokenExpiration, *newClaims)
	if err != nil {
		return newAccessToken, err
	}

	return dto.AccessToken{
		Token:     token,
		ExpiresIn: expiresIn,
	}, nil
}

func (s *AuthService) ValidateAppKey(ctx context.Context, appKey string) (isValid bool) {
	return appKey == s.appKey
}
