package models

import "time"

type Category struct {
	CategoryID   string
	CategoryName string
	ParentID     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
