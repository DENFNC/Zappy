package dao

import "github.com/jackc/pgx/v5/pgtype"

type (
	ProductCategoryDAO struct {
		ProductID  pgtype.UUID        `db:"product_id"`
		CategoryID pgtype.UUID        `db:"category__id"`
		AssignedAt pgtype.Timestamptz `db:"assigned_id"`
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
