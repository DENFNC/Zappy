package service

import (
	"context"
	"errors"

	"log/slog"

	"github.com/DENFNC/Zappy/auth_service/internal/domain/models"
	vaulttoken "github.com/DENFNC/Zappy/auth_service/internal/pkg/authjwt"
	"github.com/DENFNC/Zappy/auth_service/internal/utils/config"
	errpkg "github.com/DENFNC/Zappy/auth_service/internal/utils/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	emptyValue = 0
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user *models.User,
	) (uint64, error)
	GetByID(
		ctx context.Context,
		id string,
	) (*models.User, error)
	GetByAuthIdentifier(
		ctx context.Context,
		identifier string,
	) (*models.User, error)
	Update(
		ctx context.Context,
		user *models.User,
	) error
	List(
		ctx context.Context,
		pageSize uint,
		pageToken string,
	) ([]*models.User, string, error)
	Delete(
		ctx context.Context,
		id string,
	) error
}

type Auth struct {
	log   *slog.Logger
	repo  UserRepository
	vault vaulttoken.VaultKMS
	cfg   *config.Config
}

func NewAuth(
	log *slog.Logger,
	repo UserRepository,
	vault vaulttoken.VaultKMS,
	cfg *config.Config,
) *Auth {
	return &Auth{
		log:   log,
		repo:  repo,
		vault: vault,
		cfg:   cfg,
	}
}

func (a *Auth) UserRegister(
	ctx context.Context,
	user *models.User,
) (string, uint64, error) {
	const op = "authservice.Auth.UserRegister"

	log := a.log.With("op", op)

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("password hash generation error", slog.String("error", err.Error()))
		return "", emptyValue, errpkg.ErrInternal
	}

	token, err := a.generateToken()
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", emptyValue, errpkg.ErrInternal
	}

	userID, err := a.repo.Create(ctx, &models.User{
		Email:    user.Email,
		Username: user.Username,
		Password: passHash,
	})
	if err != nil {
		log.Error(
			"error saving the user",
			slog.String("error", err.Error()),
		)
		var appErr *errpkg.AppError
		if errors.As(err, &appErr) && appErr.Code == errpkg.ErrUniqueViolation.Code {
			return "", emptyValue, errpkg.ErrUserAlreadyExists
		}
		return "", emptyValue, errpkg.ErrInternal
	}

	return token, userID, nil
}

func (a *Auth) UserLogin(
	ctx context.Context,
	identifier string,
	password string,
) (string, error) {
	const op = "authservice.Auth.UserLogin"

	log := a.log.With("op", op)

	if identifier == "" || password == "" {
		return "", errpkg.ErrInvalidCredentials
	}

	res, err := a.repo.GetByAuthIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, errpkg.ErrUserNotFound) {
			log.Debug(
				"user not found",
				slog.String("identifier", identifier),
			)
			return "", errpkg.ErrUserNotFound
		}
		log.Error(
			"failed to execute database query",
			slog.String("error", err.Error()),
		)
		return "", errpkg.ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword(res.Password, []byte(password)); err != nil {
		log.Debug(
			"invalid password",
			slog.String("error", err.Error()),
		)
		return "", errpkg.ErrInvalidCredentials
	}

	token, err := a.generateToken()
	if err != nil {
		log.Error(
			"failed to generate token",
			slog.String("error", err.Error()),
		)
		return "", errpkg.ErrInternal
	}

	return token, nil
}

func (a *Auth) UpdateRefreshToken(
	ctx context.Context,
	token string,
) (string, error) {
	const op = "auth.Refresh"
	log := a.log.With("op", op)

	if token == "" {
		log.Error("Empty token provided for refresh")
		return "", errpkg.ErrInvalidToken
	}

	if err := vaulttoken.Verify(token); err != nil {
		log.Error(
			"Failed to verify token",
			slog.String("error", err.Error()),
		)
		return "", errpkg.ErrInvalidToken
	}

	newToken, err := a.generateToken()
	if err != nil {
		log.Error(
			"Failed to generate new token",
			slog.String("error", err.Error()),
		)
		return "", errpkg.ErrInternal
	}

	return newToken, nil
}

func (a *Auth) ListUsers(
	ctx context.Context,
	pageSize uint,
	pageToken string,
) ([]*models.User, string, error) {
	const op = "auth.ListUsers"
	log := a.log.With("op", op)

	users, nextPageToken, err := a.repo.List(ctx, pageSize, pageToken)
	if err != nil {
		log.Error(
			"Failed to list users from repository",
			slog.String("error", err.Error()),
		)
		return nil, "", errpkg.ErrInternal
	}

	return users, nextPageToken, nil
}

func (a *Auth) generateToken() (string, error) {
	token, err := vaulttoken.Generate(
		a.vault,
		a.cfg.Vault.Issuer,
		a.cfg.Vault.KeyName,
		a.cfg.Vault.Expires,
	)
	if err != nil {
		return "", errpkg.New(
			"INTERNAL_SERVER",
			"failed to generate token", err,
		)
	}
	return token, nil
}
