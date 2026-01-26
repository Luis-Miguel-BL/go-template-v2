package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/service"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/controller"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
	telemetry telemetry.Telemetry
	cfg       *config.Config
	log       logger.Logger

	authService service.AuthService

	leadController *controller.LeadController
	authController *controller.AuthController
}

func NewServer(cfg *config.Config, log logger.Logger, telemetry telemetry.Telemetry, authService service.AuthService, leadController *controller.LeadController, authController *controller.AuthController) *Server {
	server := &Server{
		Echo:           echo.New(),
		cfg:            cfg,
		log:            log,
		telemetry:      telemetry,
		authService:    authService,
		leadController: leadController,
		authController: authController,
	}

	server.setup()
	return server
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		address := fmt.Sprintf(":%d", s.cfg.Server.Port)
		s.log.Info("Starting HTTP server on " + address)
		if err := s.Echo.Start(address); err != nil && err != http.ErrServerClosed {
			s.Echo.Logger.Fatal("Error starting the server:", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Shutting down HTTP server...")
	return s.Echo.Shutdown(ctx)
}
