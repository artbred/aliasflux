package chats

import (
	"github.com/artbred/aliasflux/src/domain/models"
	"github.com/artbred/aliasflux/src/services/api/internal"
)

type CreateChatRequest struct {
	UserID       string              `json:"user_id"  validate:"required,uuid"`
	ChatSettings models.SettingsUser `json:"settings" validate:"required"`
}

type CreateChatResponse struct {
	internal.BaseResponse
	ChatID      string  `json:"chat_id"`
	PaymentLink *string `json:"payment_link"`
}

type GetChatRequest struct {
	ChatID string `param:"id" validate:"required,uuid"`
	Offset int    `query:"offset" validate:"omitempty,gte=0,number"`
}

type GetChatResponse struct {
	internal.BaseResponse
	Chat   *models.Chat `json:"chat"`
	Offset int          `json:"offset"`
}

type WebsocketConnectRequest struct {
	ChatID string `query:"id" validate:"required,uuid"`
}

type ListChatConfigurationsResponse struct {
	internal.BaseResponse
	Configurations []models.PlatformSettings `json:"configurations"`
}
