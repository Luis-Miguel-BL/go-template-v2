package middleware

import (
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/service"
	"github.com/labstack/echo/v4"
)

func NewValidateAppKeyMiddleware(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ctx := c.Request().Context()
			appKey := c.Request().Header.Get("x-api-key")

			if !authService.ValidateAppKey(ctx, appKey) {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			return next(c)
		}
	}
}
