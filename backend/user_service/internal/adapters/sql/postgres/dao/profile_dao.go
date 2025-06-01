package dao

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ProfileDAO struct {
	ProfileID  pgtype.UUID      `db:"profile_id"`
	AuthUserID pgtype.UUID      `db:"auth_user_id"`
	FirstName  pgtype.Text      `db:"first_name"`
	LastName   pgtype.Text      `db:"last_name"`
	CreatedAt  pgtype.Timestamp `db:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at"`
}
