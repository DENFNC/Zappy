package models

import (
	"time"
)

type Product struct {
	ProductID   string
	ProductName string
	Description string
	Price       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
