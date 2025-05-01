package models

import "time"

type ProductCategory struct {
	ProductID  string
	CategoryID string
	AssignedAt time.Time
}
