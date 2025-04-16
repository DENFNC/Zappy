package models

import "time"

type Profile struct {
	ProfileID  uint32
	AuthUserID uint32
	FirstName  string
	LastName   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewProfile(
	authUserID uint32,
	firstName string,
	lastName string,
) *Profile {
	return &Profile{
		AuthUserID: authUserID,
		FirstName:  firstName,
		LastName:   lastName,
	}
}
