package dao

import "github.com/jackc/pgx/v5/pgtype"

type (
	UserDAO struct {
		ID        pgtype.Int8        `db:"id"`
		Email     string             `db:"email"`
		Username  string             `db:"username"`
		Password  []byte             `db:"password_hash"`
		CreatedAt pgtype.Timestamptz `db:"created_at"`
		UpdatedAt pgtype.Timestamptz `db:"updated_at"`
	}

	RoleDAO struct {
		ID        pgtype.Int8        `db:"id"`
		Name      string             `db:"name"`
		CreatedAt pgtype.Timestamptz `db:"created_at"`
		UpdatedAt pgtype.Timestamptz `db:"updated_at"`
	}

	UserRoleDAO struct {
		UserID     pgtype.Int8        `db:"user_id"`
		RoleID     pgtype.Int8        `db:"role_id"`
		AssignedAt pgtype.Timestamptz `db:"assigned_at"`
	}
)
