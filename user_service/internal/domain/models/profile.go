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

func NewProfile(
	authUserID string,
	firstName string,
	lastName string,
) *Profile {
	return &Profile{
		AuthUserID: authUserID,
		FirstName:  firstName,
		LastName:   lastName,
	}
}
