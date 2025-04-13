package authservice

import (
	"context"
	"errors"

	"log/slog"

	"github.com/DENFNC/Zappy/internal/config"
	"github.com/DENFNC/Zappy/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/internal/errors"
	vaulttoken "github.com/DENFNC/Zappy/internal/pkg/authjwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	emptyValue = 0
)

type Auth struct {
	log      *slog.Logger
	repo     models.UserRepository
	vault    vaulttoken.VaultKMS
	tokenCfg config.ConfigVault
}

func NewAuth(
	log *slog.Logger,
	repo models.UserRepository,
	vault vaulttoken.VaultKMS,
	tokenCfg config.ConfigVault,
) *Auth {
	return &Auth{
		log:      log,
		repo:     repo,
		vault:    vault,
		tokenCfg: tokenCfg,
	}
}

func (a *Auth) Register(
	ctx context.Context,
	username string,
	email string,
	password string,
) (string, uint64, error) {
	const op = "auth.Register"
	log := a.log.With("op", op)

	if username == "" || email == "" || password == "" {
		log.Error("empty registration field")
		return "", emptyValue, errpkg.New("INVALID_CREDENTIALS", "registration fields cannot be empty", nil)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("password hash generation error", slog.String("error", err.Error()))
		return "", emptyValue, errpkg.New("INTERNAL_SERVER", "failed to generate password hash", err)
	}

	user := models.NewUser(
		username,
		email,
		passHash,
	)

	token, err := a.generateToken()
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", emptyValue, errpkg.New("INTERNAL_SERVER", "token generation failed", err)
	}

	userID, err := a.repo.Create(ctx, user)
	if err != nil {
		log.Error("error saving the user", slog.String("error", err.Error()))
		return "", emptyValue, errpkg.New("INTERNAL_SERVER", "failed to save user", err)
	}

	return token, userID, nil
}

func (a *Auth) Login(
	ctx context.Context,
	identifier string,
	password string,
) (string, error) {
	const op = "auth.Login"
	log := a.log.With("op", op)

	if identifier == "" || password == "" {
		log.Error("empty login credentials")
		return "", errpkg.New("INVALID_CREDENTIALS", "login credentials cannot be empty", nil)
	}

	res, err := a.repo.GetByAuthIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, errpkg.ErrUserNotFound) {
			log.Debug("user not found", slog.String("identifier", identifier))
			return "", errpkg.New("INVALID_CREDENTIALS", "user not found", err)
		}
		log.Error("failed to execute database query", slog.String("error", err.Error()))
		return "", errpkg.New("INTERNAL_SERVER", "database query failure", err)
	}

	if err := bcrypt.CompareHashAndPassword(res.Password, []byte(password)); err != nil {
		log.Debug("invalid password", slog.String("error", err.Error()))
		return "", errpkg.New("INVALID_CREDENTIALS", "invalid password", err)
	}

	token, err := a.generateToken()
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", errpkg.New("INTERNAL_SERVER", "token generation failed", err)
	}

	return token, nil
}

func (a *Auth) Refresh(
	ctx context.Context,
	token string,
) (string, error) {
	const op = "auth.Refresh"
	log := a.log.With("op", op)

	if token == "" {
		log.Error("Empty token provided for refresh")
		return "", errpkg.New("INVALID_TOKEN", "token is empty", nil)
	}

	if err := vaulttoken.Verify(token); err != nil {
		log.Error("Failed to verify token", slog.String("error", err.Error()))
		return "", errpkg.New("INVALID_TOKEN", "token verification failed", err)
	}

	newToken, err := a.generateToken()
	if err != nil {
		log.Error("Failed to generate new token", slog.String("error", err.Error()))
		return "", errpkg.New("INTERNAL_SERVER", "token generation failed", err)
	}

	return newToken, nil
}

func (a *Auth) generateToken() (string, error) {
	token, err := vaulttoken.Generate(
		a.vault,
		a.tokenCfg.Issuer,
		a.tokenCfg.KeyName,
		a.tokenCfg.Expires,
	)
	if err != nil {
		return "", errpkg.New("INTERNAL_SERVER", "failed to generate token", err)
	}
	return token, nil
}
