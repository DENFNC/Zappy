package models

import "time"

type ProductImage struct {
	ImageID   string
	ProductID string
	URL       string
	ALT       string
	ObjectKey string
	CreatedAt time.Time
	UpdatedAt time.Time
}
