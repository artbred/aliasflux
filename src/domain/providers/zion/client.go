package zion

import (
	"context"
	"github.com/artbred/aliasflux/src/domain/models"
	"github.com/artbred/aliasflux/src/pkg/storages/natscli"
)

const (
	CreateBusinessNamesSubject natscli.Subject = "zion.create_business_names"
)

type CreateBusinessNamesRequest struct {
	Messages       []models.ChatMessage `json:"all_messages"`
	PreferredNames []string             `json:"preferred_names"`
	Subject        natscli.Subject      `json:"-"`
	Message        string               `json:"message"`
}

type CreateBusinessNamesResponse struct {
	Names   []string `json:"names"`
	Message string   `json:"message"`
}

func (r *CreateBusinessNamesRequest) Send(ctx context.Context) (res CreateBusinessNamesResponse, err error) {
	err = natscli.PerformRequest(ctx, r.Subject, r, &res)
	return res, err
}

func NewCreateBusinessName(chat *models.Chat, message string) *CreateBusinessNamesRequest {
	return &CreateBusinessNamesRequest{
		Subject:        CreateBusinessNamesSubject,
		Messages:       chat.Messages,
		PreferredNames: chat.Settings.PreferredNames,
		Message:        message,
	}
}
