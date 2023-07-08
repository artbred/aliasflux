package redisdb

import "fmt"

type Key string

const (
	SettingsKey Key = "settings"
)

func (k Key) String() string {
	return string(k)
}

func (k Key) Build(id string) string {
	return fmt.Sprintf("%s:%s", k, id)
}
