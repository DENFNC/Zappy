package repo

import (
	"context"

	"github.com/DENFNC/Zappy/internal/domain/models"
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
				return err
			}

			return nil
		})
	})
	return userID, err
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
		return nil, err
	}

	return &user, nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	panic("implement me!")
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
