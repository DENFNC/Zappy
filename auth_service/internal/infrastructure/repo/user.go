package repo

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/internal/errors"
	psql "github.com/DENFNC/Zappy/internal/storage/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	*psql.Storage
}

func NewUser(DB *psql.Storage) *User {
	return &User{
		Storage: DB,
	}
}

func (u *User) Create(ctx context.Context, user *models.User) (uint64, error) {
	var userID uint64
	err := u.DB.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		return u.WithTx(ctx, c, func(tx pgx.Tx) error {
			if err := tx.QueryRow(ctx, `
				INSERT INTO users(email, username, password_hash)
				VALUES($1, $2, $3)
				RETURNING id`,
				user.Email, user.Username, user.Password,
			).Scan(&userID); err != nil {
				return errpkg.New("CREATE_USER_ERROR", "failed to create user", err)
			}

			return nil
		})
	})
	if err != nil {
		return 0, errpkg.New("CREATE_USER_TX_ERROR", "transaction error while creating user", err)
	}
	return userID, nil
}

func (u *User) GetByID(ctx context.Context, id string) (*models.User, error) {
	row := u.DB.QueryRow(
		ctx,
		`SELECT
			id,
			email,
			username,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE id = $1`,
		id,
	)

	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, errpkg.New("GET_USER_BY_ID_ERROR", "failed to get user by id", err)
	}

	return &user, nil
}

func (u *User) GetByAuthIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	row := u.DB.QueryRow(
		ctx,
		`SELECT
			id,
			email,
			username,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE username = $1 OR email = $1
		LIMIT 1`,
		identifier,
	)

	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrUserNotFound
		}
		return nil, errpkg.New("GET_USER_BY_AUTH_ERROR", "failed to scan user", err)
	}

	return &user, nil
}

func (u *User) Update(ctx context.Context, user *models.User) error {
	panic("implement me!")
}

func (u *User) Delete(ctx context.Context, id string) error {
	panic("implement me!")
}

func (u *User) List(ctx context.Context) ([]*models.User, error) {
	panic("implement me!")
}
