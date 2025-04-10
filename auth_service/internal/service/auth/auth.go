package auth

import "context"

type UserSaver interface {
	UserSave(
		ctx context.Context,
		username string,
		email string,
		passHash []byte,
	) (uid string, err error)
}
