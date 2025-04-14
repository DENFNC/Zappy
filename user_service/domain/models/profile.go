package models

import "time"

type Profile struct {
	ID        uint64
	FirstName string
	LastName  string
	AvatarURL []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProfile(
	firstName, lastName string,
	avatarURL []string,
) *Profile {
	return &Profile{
		FirstName: firstName,
		LastName:  lastName,
		AvatarURL: avatarURL,
	}
}
