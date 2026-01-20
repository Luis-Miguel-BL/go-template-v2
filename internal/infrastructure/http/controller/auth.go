package controller

import (
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/service"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/payload"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService}
}

func (c *AuthController) Authorization(ctx echo.Context) error {
	context := ctx.Request().Context()

	accessToken, err := c.authService.GenerateToken(context)
	if err != nil {
		return Error(ctx, err)
	}

	return Ok(ctx, payload.AuthorizationResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   accessToken.ExpiresIn,
	})
}
