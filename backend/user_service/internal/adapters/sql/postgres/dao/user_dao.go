package dao

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserDAO struct {
	UserID    pgtype.UUID      `db:"user_id"`
	Email     pgtype.Text      `db:"email"`
	Password  pgtype.Text      `db:"password"`
	CreatedAt pgtype.Timestamp `db:"created_at"`
	UpdatedAt pgtype.Timestamp `db:"updated_at"`
}
