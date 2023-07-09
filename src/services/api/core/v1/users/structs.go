package users

import "github.com/artbred/aliasflux/src/services/api/internal"

type CreateUserResponse struct {
	internal.BaseResponse
	UserID string `json:"user_id"`
}
