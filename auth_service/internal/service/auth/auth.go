package authservice

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	emptyValue uint64 = 0
)

type Auth struct {
	log  *slog.Logger
	repo models.UserRepository
}

func NewAuth(log *slog.Logger, repo models.UserRepository) *Auth {
	return &Auth{
		log:  log,
		repo: repo,
	}
}

func (a *Auth) Register(
	ctx context.Context,
	username string,
	email string,
	password string,
) (string, uint64, error) {
	const op = "auth.UserSave"

	log := a.log.With("op", op)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Debug(
			"password hash generation error",
			slog.Any("error", err),
		)
		return "", emptyValue, err
	}

	user := models.NewUser(
		username,
		email,
		passHash,
	)

	userID, err := a.repo.Create(ctx, user)
	if err != nil {
		log.Info(
			"error saving the user",
			slog.Any("error", err),
		)
		return "", emptyValue, err
	}
	log.Info(
		"A new user has been registered",
		"user_id", userID,
		"username", username,
	)

	return "", userID, nil
}

func (a *Auth) Login(
	ctx context.Context,
	authType string,
	password string,
) (string, error) {
	panic("TODO")
}

func (a *Auth) Refresh(
	ctx context.Context,
	token string,
) (string, error) {
	panic("TODO")
}
