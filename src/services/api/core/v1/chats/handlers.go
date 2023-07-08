package chats

import (
	"github.com/artbred/aliasflux/src/domain/flux"
	"github.com/artbred/aliasflux/src/domain/models"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/services/api/internal"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"reflect"
)

// CreateChatHandler
// @Summary      Create chat
// @Tags         chats
// @Produce      json
// @Accept 		 json
// @Param        data body CreateChatRequest true "Chat"
// @Success      200  {object} CreateChatResponse
// @Router       /chats/create [post]
func (r *Router) CreateChatHandler(c echo.Context) error {
	req := &CreateChatRequest{}

	if err := internal.ValidateRequest(c, req); err != nil {
		return err
	}

	t := reflect.TypeOf(req.ChatConfig)
	v := reflect.ValueOf(req.ChatConfig)

	settings := map[flux.SettingKey]interface{}{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		settingKey := field.Tag.Get("setting")

		if settingKey == "" {
			common.Logger.Errorf("failed to get setting key from tag")
			return internal.InternalServerErrorResponse(c)
		}

		fieldValue := v.Field(i).Interface()
		settings[flux.SettingKey(settingKey)] = fieldValue
	}

	isValid, err := models.CheckSettingsAreValid(settings)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to check settings are valid")
		return internal.InternalServerErrorResponse(c)
	}

	logrus.Print(isValid)

	return nil

	//valid, err := models.CheckSettingIsValid(flux.SettingKeyPlatform, req.ChatConfig.Platform)
	//if err != nil {
	//	common.Logger.WithError(err).Errorf("failed to check setting is valid")
	//	return internal.InternalServerErrorResponse(c)
	//}
	//
	//logrus.Print(valid)
	//
	//return nil
}
