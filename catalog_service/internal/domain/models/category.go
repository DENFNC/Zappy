package models

import (
	"time"
)

type Category struct {
	CategoryID   string    `db:"category_id"`
	CategoryName string    `db:"category_name"`
	ParentID     *string   `db:"parent_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
