package dao

import "github.com/jackc/pgx/v5/pgtype"

type ReviewDAO struct {
	ReviewID  pgtype.UUID      `db:"review_id"`
	ProductId pgtype.UUID      `db:"product_id"`
	ProfileID pgtype.UUID      `db:"profile_id"`
	Rating    uint32           `db:"rating"`
	Comment   pgtype.Text      `db:"comment"`
	CreatedAt pgtype.Timestamp `db:"created_at"`
	UpdatedAt pgtype.Timestamp `db:"updated_at"`
}
