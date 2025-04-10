package models

import (
	"context"
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

type UserRepository interface {
	Create(ctx context.Context, user *User) (uint64, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*User, error)
}

func NewUser(
	email string,
	username string,
	password []byte,
) *User {
	return &User{
		Email:    email,
		Username: username,
		Password: password,
	}
}
