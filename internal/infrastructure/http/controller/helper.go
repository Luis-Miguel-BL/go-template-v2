package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Ok(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, data)
}

func Error(ctx echo.Context, data any) error {
	message := "Internal Server Error"
	code := http.StatusInternalServerError
	if _, ok := data.(error); ok {
		message = data.(error).Error()
	}
	return ctx.JSON(code, map[string]any{"message": message})
}
