package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	ProductID   string
	ProductName string
	Description string
	Price       pgtype.Numeric
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
