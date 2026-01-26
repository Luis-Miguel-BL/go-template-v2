package auth

import (
	"context"
	"time"
)

type JWTHelper[Claims any] interface {
	GenerateToken(ctx context.Context, tokenExpiration time.Duration, claims Claims) (token string, expiresIn int64, err error)
	ValidateToken(encodedToken string) (claims *Claims, err error)
}

const (
	tokenClaimsContextKey = "token"
)

func WithContext[Claims any](ctx context.Context, token *Claims) (newContext context.Context) {
	newContext = context.WithValue(ctx, tokenClaimsContextKey, token)

	return newContext
}

func FromContext[Claims any](ctx context.Context) (tokenClaims *Claims) {
	t := ctx.Value(tokenClaimsContextKey)
	if t == nil {
		return nil
	}

	tokenClaims, ok := t.(*Claims)
	if !ok {
		return nil
	}

	return tokenClaims
}
