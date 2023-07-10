package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artbred/aliasflux/src/pkg/storages/postgres"
	"time"
)

type Chat struct {
	ID       string   `db:"id" json:"id"`
	Platform Platform `db:"platform" json:"platform"`
	UserID   string   `db:"user_id" json:"user_id"`

	ChatFeatures ChatFeatures `db:"features" json:"chat_features"`
	Settings     SettingsUser `db:"settings" json:"settings"`

	CreatedAt time.Time  `db:"created_at" json:"-"`
	DeletedAt *time.Time `db:"deleted_at"  json:"-"`

	Messages []ChatMessage `json:"messages"`

	UserMessagesCount int16 `db:"user_messages_count" json:"user_messages_count"`
}

type ChatMessage struct {
	ID        int64      `db:"id" json:"-"`
	ChatID    string     `db:"chat_id" json:"-"`
	Message   string     `db:"message" json:"message"`
	Sender    Sender     `db:"sender" json:"sender"`
	CreatedAt time.Time  `db:"created_at" json:"-"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

func (c *Chat) Create(ctx context.Context) (err error) {
	conn := postgres.Connection()

	query := `
		INSERT INTO chats (id, platform, user_id, settings, created_at, feature_user_messages) VALUES ($1, $2, $3, $4::jsonb, now(), $5);
	`

	_, err = conn.ExecContext(ctx, query, c.ID, c.Platform, c.UserID, c.Settings, c.ChatFeatures.UserMessagesLimit)
	if err != nil {
		err = fmt.Errorf("can't create chat %s: %v", c.ID, err)
		return
	}

	return
}

func GetChat(ctx context.Context, id string) (chat *Chat, err error) {
	conn := postgres.Connection()
	chat = &Chat{}

	query := `SELECT
				c.id,
				c.user_id,
				c.platform,
				c.settings,
				c.feature_user_messages "features.feature_user_messages",
				(
					SELECT COUNT(1)
					FROM chats_messages m
					WHERE m.chat_id = c.id AND m.sender = $2
				) AS user_messages_count
			FROM chats c
			WHERE c.id = $1;
	`

	err = conn.GetContext(ctx, chat, query, id, SenderUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get chat %s: %v", chat.ID, err)
	}

	query = `
		SELECT * FROM chats_messages WHERE chat_id = $1 AND deleted_at IS NULL AND sender::text NOT LIKE '%system%' ORDER BY created_at DESC;
	`

	err = conn.SelectContext(ctx, &chat.Messages, query, id)
	if err != nil {
		err = fmt.Errorf("can't get messages for chat %s: %v", chat.ID, err)
	}

	return
}

func GetMessagesForChat(ctx context.Context, chatID string) (messages []ChatMessage, err error) {
	conn := postgres.Connection()

	query := `
		SELECT * FROM chats_messages WHERE chat_id = $1 AND deleted_at IS NULL AND sender::text NOT LIKE '%system%' ORDER BY created_at DESC;
	`

	err = conn.SelectContext(ctx, &messages, query, chatID)
	if err != nil {
		err = fmt.Errorf("can't get messages for chat %s: %v", chatID, err)
	}

	return
}

func (c *Chat) LoadMessages(ctx context.Context, onlyForUser bool) (err error) {
	conn := postgres.Connection()

	query := `SELECT * FROM chats_messages WHERE chat_id = $1 AND deleted_at IS NULL AND sender::text NOT LIKE '%system%' ORDER BY created_at DESC;`

	if !onlyForUser {
		query = `SELECT * FROM chats_messages WHERE chat_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC;`
	}

	err = conn.SelectContext(ctx, &c.Messages, query)
	if err != nil {
		err = fmt.Errorf("can't get messages for chat %s: %v", c.ID, err)
	}

	return
}

func (m *ChatMessage) Create(ctx context.Context) (err error) {
	conn := postgres.Connection()

	query := `INSERT INTO chats_messages (chat_id, message, sender, created_at) VALUES ($1, $2, $3, now());`
	_, err = conn.ExecContext(ctx, query, m.ChatID, m.Message, m.Sender)
	return
}
