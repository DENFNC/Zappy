package models

import "time"

type Review struct {
	ReviewID  string
	ProductID string
	ProfileID string
	Rating    int16
	Comments  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
