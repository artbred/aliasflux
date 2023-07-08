package models

import (
	"fmt"
	"github.com/artbred/aliasflux/src/domain/flux"
	"github.com/artbred/aliasflux/src/pkg/storages/postgres"
	"time"
)

type Chat struct {
	ID        string     `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type ChatMessage struct {
	ID        int64      `db:"id"`
	ChatID    string     `db:"chat_id"`
	Message   string     `db:"message"`
	IsSystem  bool       `db:"is_system"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type ChatConfig struct {
	ID     int64       `db:"id"`
	ChatID string      `db:"chat_id"`
	Config flux.Config `db:"config"`
}

func CreateChat(chat *Chat, config *ChatConfig) (err error) {
	conn := postgres.Connection()
	tx := conn.MustBegin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.NamedExec("INSERT INTO chats (id, created_at) VALUES (:id, now());", chat)
	if err != nil {
		err = fmt.Errorf("can't add chat %s: %v", chat.ID, err)
		return
	}

	_, err = tx.NamedExec("INSERT INTO chat_configs (chat_id, config) VALUES (:chat_id, :config);", config)
	if err != nil {
		err = fmt.Errorf("can't add chat config %s: %v", chat.ID, err)
	}

	return
}
