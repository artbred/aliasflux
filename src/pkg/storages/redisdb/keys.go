package redisdb

import "fmt"

type Key string

const (
	SettingsKey      Key = "settings"
	OpenChats        Key = "open_chats"
	OpenChatMessages Key = "open_chat_messages"
)

func (k Key) String() string {
	return string(k)
}

func (k Key) Build(id any) string {
	return fmt.Sprintf("%s:%s", k, id)
}
