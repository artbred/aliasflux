package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artbred/aliasflux/src/domain/flux"
	"github.com/artbred/aliasflux/src/pkg/storages/postgres"
	"github.com/jmoiron/sqlx/types"
	"time"
)

type Chat struct {
	ID       string        `db:"id" json:"id"`
	Platform flux.Platform `db:"platform" json:"platform"`

	UserID   string         `db:"user_id" json:"user_id"`
	Settings types.JSONText `db:"settings" json:"settings"`

	FeatureUserMessages int16 `db:"feature_user_messages" json:"feature_user_messages"`

	CreatedAt time.Time  `db:"created_at" json:"-"`
	DeletedAt *time.Time `db:"deleted_at"  json:"-"`

	Messages []ChatMessage `db:"-" json:"messages"`
}

type ChatMessage struct {
	ID        int64      `db:"id" json:"-"`
	ChatID    string     `db:"chat_id" json:"-"`
	Message   string     `db:"message" json:"message"`
	IsSystem  bool       `db:"is_system" json:"is_system"`
	CreatedAt time.Time  `db:"created_at" json:"-"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

func (c *Chat) Create(ctx context.Context) (err error) {
	tx := postgres.Connection().MustBeginTx(ctx, nil)
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
		INSERT INTO chats (id, platform, user_id, settings, created_at, feature_user_messages) VALUES ($1, $2, $3, $4::jsonb, now(), $5);
	`

	_, err = tx.ExecContext(ctx, query, c.ID, c.Platform, c.UserID, c.Settings, c.FeatureUserMessages)
	if err != nil {
		err = fmt.Errorf("can't create chat %s: %v", c.ID, err)
		return
	}

	if len(c.Messages) > 0 {
		query = `
			INSERT INTO chats_messages (chat_id, message, is_system, created_at) VALUES
		`

		for i, message := range c.Messages {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("('%s', '%s', %t, now())", message.ChatID, message.Message, message.IsSystem)
		}

		_, err = tx.ExecContext(ctx, query)
	}

	if err != nil {
		err = fmt.Errorf("can't create messages for chat %s: %v", c.ID, err)
	}

	return
}

func GetChatWithMessages(ctx context.Context, id string, offset, limit int) (chat *Chat, err error) {
	conn := postgres.Connection()
	chat = &Chat{}

	query := `
		SELECT * FROM chats WHERE id = $1;
	`

	err = conn.GetContext(ctx, chat, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		err = fmt.Errorf("can't get chat %s: %v", chat.ID, err)
		return
	}

	query = `
		SELECT * FROM chats_messages WHERE chat_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC OFFSET $2 LIMIT $3;
	`

	err = conn.SelectContext(ctx, &chat.Messages, query, id, offset, limit)
	if err != nil {
		err = fmt.Errorf("can't get messages for chat %s: %v", chat.ID, err)
		return
	}

	return

}

func GetUserWithChat(ctx context.Context, id string) (user *User, chat *Chat, err error) {
	tx := postgres.Connection().MustBeginTx(ctx, nil)
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	chat = &Chat{}

	err = tx.Get(chat, "SELECT * FROM chats WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		err = fmt.Errorf("failed to get chat by id: %w", err)
		return
	}

	user = &User{}

	err = tx.Get(user, "SELECT * FROM users WHERE id = $1", chat.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		err = fmt.Errorf("failed to get user by id: %w", err)
		return
	}

	return
}
