package models

import "time"

type Role struct {
	ID        int64
	RoleName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
