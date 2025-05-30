package models

import "time"

type Review struct {
	ReviewID  string
	ProductID string
	ProfileID string
	Rating    uint32
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
