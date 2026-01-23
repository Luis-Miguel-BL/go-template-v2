package http

import (
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/middleware"
	"github.com/labstack/echo/v4"
)

func (s *Server) setup() {
	s.Echo.HideBanner = true
	telemetryMiddlewares := s.getObservabilityMiddlewares()
	s.Echo.Use(telemetryMiddlewares...)

	s.Echo.Use(middleware.NewLoggerMiddleware(s.log))
	s.Echo.Use(middleware.NewErrorHandlerMiddleware(s.log))
	s.Echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	baseRoute := s.Echo.Group(s.cfg.Server.Prefix)
	appKeyRoute := baseRoute.Group("", middleware.NewValidateAppKeyMiddleware(s.authService))
	publicRoute := baseRoute.Group("", middleware.NewValidateSessionMiddleware(s.authService, s.telemetry))
	privateRoute := publicRoute.Group("", middleware.NewValidateLeadSessionMiddleware(s.authService, s.telemetry))

	appKeyRoute.POST("/authorization", s.authController.Authorization)
	publicRoute.POST("/leads", s.leadController.Create)
	privateRoute.POST("/leads2", s.leadController.Create)
}

func (s *Server) getObservabilityMiddlewares() []echo.MiddlewareFunc {
	telemetryMiddlewares := []echo.MiddlewareFunc{}
	tm := s.telemetry.GetServerMiddlewares()

	for _, m := range tm {
		switch mw := m.(type) {
		case echo.MiddlewareFunc:
			telemetryMiddlewares = append(telemetryMiddlewares, mw)
		case func(next echo.HandlerFunc) echo.HandlerFunc:
			telemetryMiddlewares = append(telemetryMiddlewares, echo.MiddlewareFunc(mw))
		}
	}
	return telemetryMiddlewares
}
