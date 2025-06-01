package dao

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type WishlistDAO struct {
	ItemID    pgtype.UUID      `db:"item_id"`
	ProfileID pgtype.UUID      `db:"profile_id"`
	ProductID pgtype.UUID      `db:"product_id"`
	AddedAt   pgtype.Timestamp `db:"added_at"`
	IsActive  pgtype.Bool      `db:"is_active"`
}
