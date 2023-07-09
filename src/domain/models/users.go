package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artbred/aliasflux/src/pkg/storages/postgres"
	"github.com/jmoiron/sqlx/types"
	"time"
)

type User struct {
	ID           string             `db:"id" json:"id"`
	CreateParams types.NullJSONText `db:"create_params" json:"-"`
	CreatedAt    time.Time          `db:"created_at" json:"-"`
	DeletedAt    *time.Time         `db:"deleted_at" json:"-"`
}

func (u *User) Create(ctx context.Context) (err error) {
	conn := postgres.Connection()

	_, err = conn.NamedExecContext(ctx, "INSERT INTO users (id, create_params) VALUES (:id, :create_params)", u)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return
}

func GetUserByID(ctx context.Context, id string) (user *User, err error) {
	conn := postgres.Connection()
	user = &User{}

	err = conn.GetContext(ctx, user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		err = fmt.Errorf("failed to get user by id: %w", err)
	}

	return
}
