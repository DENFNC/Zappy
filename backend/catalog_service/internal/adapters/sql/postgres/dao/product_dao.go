package dao

import "github.com/jackc/pgx/v5/pgtype"

type (
	ProductCategoryDAO struct {
		ProductID  pgtype.UUID        `db:"product_id"`
		CategoryID pgtype.UUID        `db:"category__id"`
		AssignedAt pgtype.Timestamptz `db:"assigned_id"`
	}

	ProductImageDAO struct {
		ImageID   string             `db:"image_id"`
		ProductID pgtype.UUID        `db:"product__id"`
		URL       string             `db:"url"`
		ALT       string             `db:"alt"`
		CreatedAt pgtype.Timestamptz `db:"created_at"`
		UpdatedAt pgtype.Timestamptz `db:"updated_at"`
	}

	ProductDAO struct {
		ProductID   pgtype.UUID        `db:"product_id"`
		ProductName string             `db:"product_name"`
		Description string             `db:"description"`
		Price       pgtype.Numeric     `db:"price"`
		CreatedAt   pgtype.Timestamptz `db:"created_at"`
		UpdatedAt   pgtype.Timestamptz `db:"updated_at"`
	}
)
