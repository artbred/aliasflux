package internal

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type BaseResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func InternalServerErrorResponse(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, BaseResponse{
		Ok:      false,
		Message: "Please try again later",
	})
}

func ValidateRequest(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Ok:      false,
			Message: err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Ok:      false,
			Message: err.Error(),
		})
	}

	return nil
}
