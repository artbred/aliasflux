package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Platform string

const (
	PlatformDomain Platform = "domain"
)

/////////////////////////////

type Sender string

const (
	SenderUser     Sender = "user"
	SenderAI       Sender = "ai"
	SenderAISystem Sender = "ai_system"
	SenderSystem   Sender = "system"
)

/////////////////////////////

type SettingKey string

const (
	SettingTld SettingKey = "tld"
)

/////////////////////////////

type TldSettings struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

/////////////////////////////

const SettingPreferredNames SettingKey = "preferred_names"
const SettingMaxLength SettingKey = "max_length"

type SettingsUser struct {
	Platform       Platform `validate:"required,max=10" json:"platform"`
	PreferredNames []string `json:"preferred_names" validate:"max=50"`

	Tld []TldSettings `settings:"tld" json:"tld"`
}

func (s *SettingsUser) Scan(value interface{}) error {
	var data []byte

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return nil
	}

	return json.Unmarshal(data, s)
}

func (s *SettingsUser) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type ChatFeatures struct {
	UserMessagesLimit int16 `db:"feature_user_messages" json:"messages_limit"`
}

func (c *ChatFeatures) Scan(value interface{}) error {
	var data []byte

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return nil
	}

	return json.Unmarshal(data, c)
}

func (c *ChatFeatures) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// SettingKeyFreeChatFeatures for a chat by user
const SettingKeyFreeChatFeatures SettingKey = "free_chat_features"
