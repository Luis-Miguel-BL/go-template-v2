package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/auth"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/service"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/labstack/echo/v4"
)

func NewValidateSessionMiddleware(authService service.AuthService, telemetry telemetry.Telemetry) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ctx := c.Request().Context()
			accessToken, ok := getBearerAuth(c.Request())
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			tokenClaims, err := authService.ValidateToken(ctx, accessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			ctx = auth.WithContext(ctx, tokenClaims)
			ctx = enrichContextAndLogger(ctx, *tokenClaims, telemetry)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func NewValidateLeadSessionMiddleware(authService service.AuthService, telemetry telemetry.Telemetry) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ctx := c.Request().Context()

			token := auth.FromContext[service.TokenClaims](ctx)
			if token == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}
			if token.LeadID == nil || *token.LeadID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			ctx = enrichContextAndLogger(ctx, *token, telemetry)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func getBearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access_token")
	}

	return token, token != ""
}

func enrichContextAndLogger(ctx context.Context, token service.TokenClaims, telemetry telemetry.Telemetry) context.Context {
	leadID := ""
	if token.LeadID != nil {
		leadID = *token.LeadID
	}
	sessionID := ""
	if token.SessionID != nil {
		sessionID = *token.SessionID
	}

	ctx = telemetry.AddAttributes(ctx, map[string]any{
		"lead_id":    leadID,
		"session_id": sessionID,
	})

	log := logger.FromContext(ctx)
	log = log.WithFields(map[string]any{
		"lead_id":    leadID,
		"session_id": sessionID,
		"trace_id":   telemetry.TraceIDFromContext(ctx),
	})

	return logger.NewContext(ctx, log)
}
