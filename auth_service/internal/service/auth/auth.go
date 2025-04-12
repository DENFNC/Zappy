package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DENFNC/Zappy/internal/config"
	"github.com/DENFNC/Zappy/internal/domain/models"
	"github.com/DENFNC/Zappy/internal/infrastructure/repo"
	vaulttoken "github.com/DENFNC/Zappy/internal/pkg/authjwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	emptyValue uint64 = 0
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
		return "", emptyValue, fmt.Errorf("%w: registration fields cannot be empty", ErrInvalidCredentials)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("password hash generation error", slog.String("error", err.Error()))
		return "", emptyValue, fmt.Errorf("%w: failed to generate password hash", ErrInternalServer)
	}

	user := models.NewUser(
		username,
		email,
		passHash,
	)

	token, err := a.generateToken()
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", emptyValue, fmt.Errorf("%w: token generation failed", ErrInternalServer)
	}

	userID, err := a.repo.Create(ctx, user)
	if err != nil {
		log.Error("error saving the user", slog.String("error", err.Error()))
		return "", emptyValue, fmt.Errorf("%w: failed to save user", ErrInternalServer)
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
		return "", fmt.Errorf("%w: login credentials cannot be empty", ErrInvalidCredentials)
	}

	res, err := a.repo.GetByAuthIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			log.Debug("user not found", slog.String("identifier", identifier))
			return "", fmt.Errorf("%w: user not found", ErrInvalidCredentials)
		}
		log.Error("failed to execute database query", slog.String("error", err.Error()))
		return "", fmt.Errorf("%w: database query failure", ErrInternalServer)
	}

	if err := bcrypt.CompareHashAndPassword(res.Password, []byte(password)); err != nil {
		log.Debug("invalid password", slog.String("error", err.Error()))
		return "", fmt.Errorf("%w: invalid password", ErrInvalidCredentials)
	}

	token, err := a.generateToken()
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%w: token generation failed", ErrInternalServer)
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
		log.Error("empty token provided for refresh")
		return "", fmt.Errorf("%w: token is empty", ErrInvalidToken)
	}

	if err := vaulttoken.Verify(token); err != nil {
		log.Error("failed to verify token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%w: token verification failed", ErrInvalidToken)
	}

	newToken, err := a.generateToken()
	if err != nil {
		log.Error("failed to generate new token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%w: token generation failed", ErrInternalServer)
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
		return "", err
	}
	return token, nil
}
