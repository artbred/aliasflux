package users

import (
	"github.com/artbred/aliasflux/src/domain/models"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/services/api/internal"
	"github.com/labstack/echo/v4"
	"github.com/twinj/uuid"
	"net/http"
)

// CreateUserHandler
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} CreateUserResponse
// @Router /users/create [get]
func (r *Router) CreateUserHandler(c echo.Context) error {
	user := models.User{
		ID: uuid.NewV4().String(),
	}

	if err := user.Create(c.Request().Context()); err != nil {
		common.Logger.WithError(err).Error("failed to create user")
		return internal.InternalServerErrorResponse(c)
	}

	return c.JSON(http.StatusCreated, CreateUserResponse{
		UserID: user.ID,
		BaseResponse: internal.BaseResponse{
			Ok:      true,
			Message: "User created",
		},
	})
}
