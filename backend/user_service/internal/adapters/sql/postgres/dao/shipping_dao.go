package dao

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ShippingDAO struct {
	AddressID  pgtype.UUID      `db:"address_id"`
	ProfileID  pgtype.UUID      `db:"profile_id"`
	Country    pgtype.Text      `db:"country"`
	City       pgtype.Text      `db:"city"`
	Street     pgtype.Text      `db:"street"`
	PostalCode pgtype.Text      `db:"postal_code"`
	IsDefault  pgtype.Bool      `db:"is_default"`
	CreatedAt  pgtype.Timestamp `db:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at"`
}
