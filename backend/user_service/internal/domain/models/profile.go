package models

import "time"

type Profile struct {
	ProfileID  string
	AuthUserID string
	FirstName  string
	LastName   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
