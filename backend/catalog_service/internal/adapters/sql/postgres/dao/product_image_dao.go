package dao

import "github.com/jackc/pgx/v5/pgtype"

type ProductImageDAO struct {
	ImageID   pgtype.UUID        `db:"image_id"`
	ProductID pgtype.UUID        `db:"product__id"`
	URL       string             `db:"url"`
	ALT       string             `db:"alt"`
	ObjectKey string             `db:"object_key"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}
