package flux

type SettingKey string

const SettingKeyPlatform SettingKey = "platform"

type PlatformSettings struct {
	Name string `json:"platform_name" validate:"required"`
}
