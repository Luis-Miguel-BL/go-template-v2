package service

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/dto"
)

type AuthService interface {
	GenerateToken(ctx context.Context) (accessToken dto.AccessToken, err error)
	ValidateToken(ctx context.Context, token string) (tokenClaims *TokenClaims, err error)
	RefreshToken(ctx context.Context, claims *TokenClaims) (newAccessToken dto.AccessToken, err error)
	ValidateAppKey(ctx context.Context, appKey string) (isValid bool)
}
