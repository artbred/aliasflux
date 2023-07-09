package flux

type Platform string
type SettingKey string

const (
	PlatformDomain Platform = "domain"
)

// Settings that can be specified for a chat by user and must be validated

const SettingTld SettingKey = "tld"

type TldSettings struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Settings that can be specified for a chat by user

const SettingPreferredNames SettingKey = "preferred_names"
const SettingMaxLength SettingKey = "max_length"

type Settings struct {
	Platform       Platform      `validate:"required,max=10" json:"platform"`
	Tld            []TldSettings `settings:"tld" json:"tld"`
	PreferredNames []string      `json:"preferred_names" validate:"max=50"`
	MaxLength      int16         `json:"max_length" validate:"number"`
}

// Settings for a free chat

const SettingKeyFreeChatFeatures SettingKey = "free_chat_features"

type FreeChatFeatures struct {
	MessagesLimit int16 `json:"messages_limit"`
}
