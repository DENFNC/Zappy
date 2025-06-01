package dao

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentDAO struct {
	PaymentID    pgtype.UUID      `db:"payment_id"`
	ProfileID    pgtype.UUID      `db:"profile_id"`
	PaymentToken pgtype.Text      `db:"payment_token"`
	IsDefault    pgtype.Bool      `db:"is_default"`
	CreatedAt    pgtype.Timestamp `db:"created_at"`
	UpdatedAt    pgtype.Timestamp `db:"updated_at"`
}
