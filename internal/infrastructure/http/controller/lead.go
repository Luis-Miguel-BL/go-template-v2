package controller

import (
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/dto"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/usecase"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/request"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/response"
	"github.com/labstack/echo/v4"
)

type LeadController struct {
	createLeadUseCase *usecase.CreateLead
}

func NewLeadController(createLeadUseCase *usecase.CreateLead) *LeadController {
	return &LeadController{createLeadUseCase}
}

func (c *LeadController) Create(ctx echo.Context) error {
	context := ctx.Request().Context()

	log := logger.FromContext(context)

	log.Debug("Creating lead")
	reqBody := request.CreateLeadRequest{}
	err := ctx.Bind(&reqBody)
	if err != nil {
		return Error(ctx, err)
	}

	output, err := c.createLeadUseCase.Execute(context, dto.CreateLeadInput(reqBody))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, response.CreateLeadResponse{
		LeadID:      output.LeadID,
		AccessToken: output.AccessToken.Token,
		ExpiresIn:   output.AccessToken.ExpiresIn,
	})
}
