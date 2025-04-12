package authservice

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidToken       = errors.New("token is invalid")
)
