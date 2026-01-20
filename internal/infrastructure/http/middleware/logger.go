package middleware

import (
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/labstack/echo/v4"
)

func NewLoggerMiddleware(log logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ctx := c.Request().Context()
			c.SetRequest(c.Request().WithContext(logger.NewContext(ctx, log)))

			start := time.Now()
			err = next(c)
			duration := time.Since(start)

			log.Debug("Completed request",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"duration", duration.String(),
				"error", err,
			)
			return err
		}
	}
}
