package models

import "time"

type ProductImage struct {
	ImageID   string
	ProductID string
	URL       string
	ALT       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
