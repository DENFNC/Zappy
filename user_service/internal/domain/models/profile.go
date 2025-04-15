package models

import "time"

type Profile struct {
	ProfileID  int
	AuthUserID int
	FirstName  string
	LastName   string
	Phone      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
