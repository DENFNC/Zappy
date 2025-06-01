package models

import (
	"time"
)

type User struct {
	ID        uint64
	Email     string
	Username  string
	Password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
