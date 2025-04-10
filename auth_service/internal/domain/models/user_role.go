package models

import "time"

type UserRole struct {
	UserID    int64
	RoleID    int64
	CreatedAt time.Time
}
