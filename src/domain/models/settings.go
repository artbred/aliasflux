package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/artbred/aliasflux/src/domain/flux"
	"github.com/artbred/aliasflux/src/pkg/storages/postgres"
	"reflect"
)

type Settings struct {
	ID       int            `db:"id" json:"-"`
	Platform *flux.Platform `db:"platform" json:"platform"`
	Key      string         `db:"key" json:"key"`
	Value    string         `db:"value" json:"value"`
}

type ChatSettings struct {
	Platform flux.Platform   `db:"platform" json:"platform"`
	Settings json.RawMessage `db:"settings" json:"settings"`
}

func GetAllSettings() (settings []Settings, err error) {
	conn := postgres.Connection()

	err = conn.Select(&settings, "SELECT * FROM settings;")
	if err != nil {
		err = fmt.Errorf("can't get all settings: %v", err)
	}

	return
}

func ValidateSettings(ctx context.Context, settings flux.Settings) (validationErr, err error) {
	t := reflect.TypeOf(settings)
	v := reflect.ValueOf(settings)

	settingsMap := map[string]interface{}{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		settingKey := field.Tag.Get("settings")

		if len(settingKey) > 0 {
			settingsMap[settingKey] = v.Field(i).Interface()
		}
	}

	conn := postgres.Connection()
	isValid := true

	for key, val := range settingsMap {
		query := "SELECT EXISTS (SELECT 1 FROM settings WHERE platform=$1 AND key=$2 AND value @> $3::jsonb);"

		err = conn.GetContext(ctx, &isValid, query, settings.Platform, key, val)
		if err != nil {
			return nil, fmt.Errorf("can't check settings: %v", err)
		}

		if !isValid {
			return fmt.Errorf("You provider invalid value for %s settings", key), nil
		}
	}

	return nil, nil
}

func GetFreeChatFeatures(ctx context.Context) (features *flux.FreeChatFeatures, err error) {
	conn := postgres.Connection()
	var val []byte

	err = conn.GetContext(ctx, &val, "SELECT value FROM settings WHERE key=$1;", flux.SettingKeyFreeChatFeatures)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		err = fmt.Errorf("can't get free chat features: %v", err)
	}

	features = &flux.FreeChatFeatures{}
	err = json.Unmarshal(val, features)

	return
}

func ListAvailableChatSettings(ctx context.Context) (settings []ChatSettings, err error) {
	conn := postgres.Connection()

	query := `
		SELECT
			platform,
			jsonb_object_agg(key, value) as settings
		FROM
			settings
		WHERE platform IS NOT NULL
		GROUP BY
			platform;
	`

	err = conn.SelectContext(ctx, &settings, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		err = fmt.Errorf("can't list available chat configurations: %v", err)
	}

	return
}
