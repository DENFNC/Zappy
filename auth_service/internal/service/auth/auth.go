package authservice

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/internal/domain/models"
	vaulttoken "github.com/DENFNC/Zappy/internal/pkg/authjwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	emptyValue uint64 = 0
)

type Auth struct {
	log   *slog.Logger
	repo  models.UserRepository
	vault vaulttoken.VaultKMS
}

func NewAuth(
	log *slog.Logger,
	repo models.UserRepository,
	vault vaulttoken.VaultKMS,
) *Auth {
	return &Auth{
		log:   log,
		repo:  repo,
		vault: vault,
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

	token, err := vaulttoken.Generate(a.vault, "test", "rsa", 10)
	if err != nil {
		log.Error(
			"failed to create token",
			slog.Any("error", err),
		)
	}

	userID, err := a.repo.Create(ctx, user)
	if err != nil {
		log.Error(
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

	return token, userID, nil
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
