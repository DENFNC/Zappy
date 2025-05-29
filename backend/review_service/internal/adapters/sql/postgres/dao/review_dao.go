package dao

import "github.com/jackc/pgx/v5/pgtype"

type ReviewDAO struct {
	ReviewID  pgtype.UUID      `db:"review_id"`
	ProductId pgtype.UUID      `db:"product_id"`
	ProfileID pgtype.UUID      `db:"profile_id"`
	Rating    int16            `db:"rating"`
	Comment   pgtype.Text      `db:"comments"`
	CreatedAt pgtype.Timestamp `db:"created_at"`
	UpdatedAt pgtype.Timestamp `db:"updated_at"`
}
