package models

import "time"

type Profile struct {
	ProfileID  uint64
	AuthUserID uint64
	FirstName  string
	LastName   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewProfile(
	authUserID uint64,
	firstName string,
	lastName string,
) *Profile {
	return &Profile{
		AuthUserID: authUserID,
		FirstName:  firstName,
		LastName:   lastName,
	}
}
