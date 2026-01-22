package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func NewErrorHandlerMiddleware(log logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				status, response := mapError(err)

				if status == http.StatusInternalServerError {
					log.Error("internal server error", "error", err)
				}

				if c.Response().Committed {
					return err
				}

				if c.Request().Method == http.MethodHead {
					return c.NoContent(status)
				}

				return c.JSON(status, response)
			}
			return nil
		}
	}
}

func mapError(err error) (int, ErrorResponse) {
	var (
		domainErr domain.DomainError
		echoErr   *echo.HTTPError
	)

	switch {
	case errors.As(err, &domainErr):
		return mapDomainError(domainErr)

	case errors.As(err, &echoErr):
		return mapEchoError(echoErr)

	default:
		return http.StatusInternalServerError, ErrorResponse{
			Message: "Internal Server Error",
		}
	}
}

func mapDomainError(err domain.DomainError) (int, ErrorResponse) {
	status := http.StatusInternalServerError

	switch err.Code {
	case domain.ErrCodeEntityNotFound:
		status = http.StatusNotFound
	case domain.ErrCodeInvalidInput:
		status = http.StatusBadRequest
	}

	return status, ErrorResponse{
		Message: err.Error(),
		Code:    string(err.Code),
	}
}

func mapEchoError(err *echo.HTTPError) (int, ErrorResponse) {
	msg, ok := err.Message.(string)
	if !ok {
		msg = fmt.Sprintf("%v", err.Message)
	}

	return err.Code, ErrorResponse{
		Message: msg,
	}
}
