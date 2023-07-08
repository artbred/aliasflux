package models

import (
	"encoding/json"
	"fmt"
	"github.com/artbred/aliasflux/src/domain/flux"
	"github.com/artbred/aliasflux/src/pkg/storages/postgres"
	"strings"
)

type Settings struct {
	Key   string `db:"key" json:"key"`
	Value string `db:"value" json:"value"`
}

func GetAllSettings() (settings []Settings, err error) {
	conn := postgres.Connection()

	err = conn.Select(&settings, "SELECT * FROM settings;")
	if err != nil {
		err = fmt.Errorf("can't get all settings: %v", err)
	}

	return
}

func CheckSettingsAreValid(settings map[flux.SettingKey]interface{}) (isValid bool, err error) {
	conn := postgres.Connection()

	var queries []string

	for key, val := range settings {
		valBytes, _ := json.Marshal(val)
		queries = append(queries, fmt.Sprintf("('%s'::jsonb) ['%s'] IS NOT NULL", string(valBytes), key))
	}

	query := "SELECT " + strings.Join(queries, " AND ") + " FROM settings;"

	err = conn.Get(&isValid, query)
	if err != nil {
		err = fmt.Errorf("can't check settings: %v", err)
	}

	return
}
