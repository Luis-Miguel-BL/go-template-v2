package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/auth"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type jwtHelper[Claims any] struct {
	secretKey string
}

type jwtClaims[Claims any] struct {
	Claims Claims
	jwt.RegisteredClaims
}

func NewJWTHelper[Claims any](cfg *config.Config) auth.JWTHelper[Claims] {
	return &jwtHelper[Claims]{cfg.Server.JWTSecret}
}

var invalidAccessTokenErr = errors.New("invalid access token")

func (p *jwtHelper[Claims]) GenerateToken(ctx context.Context, tokenExpiration time.Duration, claims Claims) (encodedToken string, expiresIn int64, err error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(tokenExpiration)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims[Claims]{
		Claims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})

	token, err := t.SignedString([]byte(p.secretKey))
	if err != nil {
		return encodedToken, expiresIn, err
	}

	return token, int64(tokenExpiration.Seconds()), nil
}

func (p *jwtHelper[Claims]) ValidateToken(encodedToken string) (claims *Claims, err error) {
	token, err := jwt.ParseWithClaims(string(encodedToken), &jwtClaims[Claims]{}, p.keyFunc)
	if err != nil {
		return claims, invalidAccessTokenErr
	}

	if !token.Valid {
		return claims, invalidAccessTokenErr
	}

	t, ok := token.Claims.(*jwtClaims[Claims])
	if !ok {
		return claims, invalidAccessTokenErr
	}

	if t.ExpiresAt.Before(time.Now()) {
		return claims, invalidAccessTokenErr
	}

	return &t.Claims, nil
}

func (p *jwtHelper[Claims]) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, jwt.ErrSignatureInvalid
	}
	return []byte(p.secretKey), nil
}
