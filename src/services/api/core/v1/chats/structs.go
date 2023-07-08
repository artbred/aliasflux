package chats

import (
	"github.com/artbred/aliasflux/src/domain/flux"
	"github.com/artbred/aliasflux/src/services/api/internal"
)

type CreateChatRequest struct {
	ChatConfig flux.Config `json:"chat_config" validate:"required"`
}

type CreateChatResponse struct {
	internal.BaseResponse
	ChatID string `json:"chat_id"`
}
