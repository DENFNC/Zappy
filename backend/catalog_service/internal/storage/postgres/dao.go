package psql

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CategoryDAO struct {
	CategoryID   pgtype.UUID        `db:"category_id"`
	CategoryName string             `db:"category_name"`
	ParentID     pgtype.UUID        `db:"parent_id"`
	CreatedAt    pgtype.Timestamptz `db:"created_at"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at"`
}
