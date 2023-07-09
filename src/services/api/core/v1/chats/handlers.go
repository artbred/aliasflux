package chats

import (
	"encoding/json"
	"github.com/artbred/aliasflux/src/domain/models"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/services/api/internal"
	"github.com/jmoiron/sqlx/types"
	"github.com/labstack/echo/v4"
	"github.com/twinj/uuid"
	"net/http"
)

// CreateChatHandler
// @Summary      Create chat
// @Tags         chats
// @Produce      json
// @Accept 		 json
// @Param        data body CreateChatRequest true "Chat"
// @Success      201  {object} CreateChatResponse
// @Router       /chats/create [post]
func (r *Router) CreateChatHandler(c echo.Context) error {
	req := &CreateChatRequest{}

	if err := internal.ValidateRequest(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: err.Error(),
		})
	}

	ctx := c.Request().Context()

	validationErr, err := models.ValidateSettings(ctx, req.ChatSettings)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to check settings are valid")
		return internal.InternalServerErrorResponse(c)
	}

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: validationErr.Error(),
		})
	}

	user, err := models.GetUserByID(ctx, req.UserID)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to get user by id")
		return internal.InternalServerErrorResponse(c)
	}

	if user == nil || user.DeletedAt != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: "We couldn't find a user with this id",
		})
	}

	freeChatFeatures, err := models.GetFreeChatFeatures(ctx)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to get free chat features")
		return internal.InternalServerErrorResponse(c)
	}

	if freeChatFeatures == nil || freeChatFeatures.MessagesLimit == 0 {
		// TODO: create payment link
		return nil
	}

	settingsBytes, _ := json.Marshal(req.ChatSettings)
	chat := &models.Chat{
		ID:                  uuid.NewV4().String(),
		UserID:              user.ID,
		Settings:            types.JSONText(settingsBytes),
		Platform:            req.ChatSettings.Platform,
		FeatureUserMessages: freeChatFeatures.MessagesLimit,
	}

	err = chat.Create(ctx)
	if err != nil {
		common.Logger.WithError(err).WithField("user_id", user.ID).Error("failed to create chat for user")
		return internal.InternalServerErrorResponse(c)
	}

	return c.JSON(http.StatusCreated, CreateChatResponse{
		ChatID: chat.ID,
		BaseResponse: internal.BaseResponse{
			Ok:      true,
			Message: "Chat created successfully!",
		},
	})
}

// GetChatHandler
// @Summary      Get chat
// @Tags         chats
// @Produce      json
// @Accept 		 json
// @Param        id path string true "Chat ID"
// @Param        offset query int true "Offset"
// @Success      200  {object} CreateChatResponse
// @Router       /chats/{id} [get]
func (r *Router) GetChatHandler(c echo.Context) error {
	req := &GetChatRequest{}

	if err := internal.ValidateRequest(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: err.Error(),
		})
	}

	chat, err := models.GetChatWithMessages(c.Request().Context(), req.ChatID, req.Offset, 25)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to get chat")
		return internal.InternalServerErrorResponse(c)
	}

	return c.JSON(http.StatusOK, GetChatResponse{
		Chat:   chat,
		Offset: req.Offset + len(chat.Messages),
		BaseResponse: internal.BaseResponse{
			Ok:      true,
			Message: "Chat retrieved successfully!",
		},
	})
}

func (r *Router) WebsocketHandler(c echo.Context) error {
	req := &WebsocketConnectRequest{}

	err := internal.ValidateRequest(c, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: err.Error(),
		})
	}

	user, chat, err := models.GetUserWithChat(c.Request().Context(), req.ChatID)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to get user with chat")
		return internal.InternalServerErrorResponse(c)
	}

	if user == nil || chat == nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: "We couldn't find a user or chat with this id",
		})
	}

	if user.DeletedAt != nil || chat.DeletedAt != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: "Your account or chat was deleted",
		})
	}

	return nil
}

// ListAvailableChatSettingsHandler
// @Summary      Get available chat configurations
// @Tags         chats
// @Produce      json
// @Success      200  {object} CreateChatResponse
// @Router       /chats/settings [get]
func (r *Router) ListAvailableChatSettingsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	configs, err := models.ListAvailableChatSettings(ctx)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to list available chat settings")
		return internal.InternalServerErrorResponse(c)
	}

	if configs == nil {
		common.Logger.Errorf("no available chat settings found")
		return internal.InternalServerErrorResponse(c)
	}

	return c.JSON(http.StatusOK, ListChatConfigurationsResponse{
		Configurations: configs,
		BaseResponse: internal.BaseResponse{
			Ok:      true,
			Message: "Chat configurations retrieved successfully!",
		},
	})
}
