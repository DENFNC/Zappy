package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/DENFNC/Zappy/auth_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/auth_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/auth_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/auth_service/internal/pkg/errors"
	"github.com/DENFNC/Zappy/auth_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/auth_service/internal/utils/dbutils"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// poolDAO хранит уже созданные dao.UserDAO, чтобы не аллоцировать каждый раз
var poolDAO = sync.Pool{
	New: func() any {
		return new(dao.UserDAO)
	},
}

// poolModel хранит экземпляры models.User, чтобы переиспользовать
var poolModel = sync.Pool{
	New: func() any {
		return new(models.User)
	},
}

// poolSliceUsers даёт готовый срез []*models.User (под initial capacity)
var poolSliceUsers = sync.Pool{
	New: func() any {
		slice := make([]*models.User, 0, 250)
		return &slice
	},
}

type User struct {
	*postgres.Storage
	goqu.DialectWrapper
	paginator *paginate.Paginator[dao.UserDAO]
}

func NewUser(db *postgres.Storage, coder paginate.TokenCoder) *User {
	paginator, err := paginate.NewPaginator[dao.UserDAO](db.Client, db.Dialect, coder)
	if err != nil {
		panic(err)
	}
	return &User{
		Storage:        db,
		DialectWrapper: db.Dialect,
		paginator:      paginator,
	}
}

func (repo *User) Create(
	ctx context.Context,
	user *models.User,
) (uint64, error) {
	var userID uint64
	err := repo.Client.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		return repo.WithTx(ctx, func(tx pgx.Tx) error {
			daoUser := poolDAO.Get().(*dao.UserDAO)
			daoUser.Reset()
			daoUser.Email = user.Email
			daoUser.Username = user.Username
			daoUser.Password = user.Password

			stmt, args, err := repo.Dialect.Insert("users").Rows(
				goqu.Record{
					"email":         daoUser.Email,
					"username":      daoUser.Username,
					"password_hash": daoUser.Password,
				},
			).Returning("id").Prepared(true).ToSQL()
			if err != nil {
				// Вернём dao в пул, чтобы не терять
				poolDAO.Put(daoUser)
				return errpkg.New("CREATE_USER_ERROR", "failed to build SQL query", err)
			}

			if err := tx.QueryRow(ctx, stmt, args...).Scan(&userID); err != nil {
				poolDAO.Put(daoUser)
				return errpkg.New("CREATE_USER_ERROR", "failed to create user", err)
			}

			// Возвращаем daoUser в пул
			poolDAO.Put(daoUser)
			return nil
		})
	})
	if err != nil {
		return 0, errpkg.New("CREATE_USER_TX_ERROR", "transaction error while creating user", err)
	}
	return userID, nil
}

func (repo *User) GetByID(
	ctx context.Context,
	id string,
) (*models.User, error) {
	stmt, args, err := repo.Dialect.Select(
		"id", "email", "username", "password_hash", "created_at", "updated_at",
	).From("users").Where(goqu.C("id").Eq(id)).Prepared(true).ToSQL()
	if err != nil {
		return nil, errpkg.New("GET_USER_BY_ID_ERROR", "failed to build SQL query", err)
	}

	// Берём dao из пула
	daoUser := poolDAO.Get().(*dao.UserDAO)
	daoUser.Reset()

	row := repo.Client.QueryRow(ctx, stmt, args...)
	if err := dbutils.ScanStruct(row, daoUser); err != nil {
		poolDAO.Put(daoUser)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrUserNotFound
		}
		return nil, errpkg.New("GET_USER_BY_ID_ERROR", "failed to get user by id", err)
	}

	// Берём модель из пула
	modelUser := poolModel.Get().(*models.User)
	*modelUser = models.User{
		ID:        uint64(daoUser.ID.Int64),
		Email:     daoUser.Email,
		Username:  daoUser.Username,
		Password:  daoUser.Password,
		CreatedAt: daoUser.CreatedAt.Time,
		UpdatedAt: daoUser.UpdatedAt.Time,
	}
	// Возвращаем dao обратно в пул
	poolDAO.Put(daoUser)
	return modelUser, nil
}

func (repo *User) GetByAuthIdentifier(
	ctx context.Context,
	identifier string,
) (*models.User, error) {
	stmt, args, err := repo.Dialect.Select(
		"id", "email", "username", "password_hash", "created_at", "updated_at",
	).From("users").
		Where(
			goqu.Or(
				goqu.C("username").Eq(identifier),
				goqu.C("email").Eq(identifier),
			),
		).
		Limit(1).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, errpkg.New("GET_USER_BY_AUTH_ERROR", "failed to build SQL query", err)
	}

	daoUser := poolDAO.Get().(*dao.UserDAO)
	daoUser.Reset()

	row := repo.Client.QueryRow(ctx, stmt, args...)
	if err := dbutils.ScanStruct(row, daoUser); err != nil {
		poolDAO.Put(daoUser)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrUserNotFound
		}
		return nil, errpkg.New("GET_USER_BY_AUTH_ERROR", "failed to scan user", err)
	}

	modelUser := poolModel.Get().(*models.User)
	*modelUser = models.User{
		ID:        uint64(daoUser.ID.Int64),
		Email:     daoUser.Email,
		Username:  daoUser.Username,
		Password:  daoUser.Password,
		CreatedAt: daoUser.CreatedAt.Time,
		UpdatedAt: daoUser.UpdatedAt.Time,
	}
	poolDAO.Put(daoUser)
	return modelUser, nil
}

func (repo *User) Update(
	ctx context.Context,
	user *models.User,
) error {
	const op = "repository.User.Update"

	err := repo.Client.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		return repo.WithTx(ctx, func(tx pgx.Tx) error {
			daoUser := poolDAO.Get().(*dao.UserDAO)
			*daoUser = dao.UserDAO{
				ID:       pgtype.Int8{Int64: int64(user.ID), Valid: true},
				Email:    user.Email,
				Username: user.Username,
				Password: user.Password,
			}

			stmt, args, err := repo.Dialect.Update("users").Set(
				goqu.Record{
					"email":         daoUser.Email,
					"username":      daoUser.Username,
					"password_hash": daoUser.Password,
					"updated_at":    goqu.L("NOW()"),
				},
			).Where(goqu.C("id").Eq(daoUser.ID)).Prepared(true).ToSQL()
			if err != nil {
				poolDAO.Put(daoUser)
				return errpkg.New("UPDATE_USER_ERROR", "failed to build SQL query", err)
			}

			cmd, err := tx.Exec(ctx, stmt, args...)
			if err != nil {
				poolDAO.Put(daoUser)
				return errpkg.New("UPDATE_USER_ERROR", "failed to update user", err)
			}

			if cmd.RowsAffected() == 0 {
				poolDAO.Put(daoUser)
				return errpkg.ErrUserNotFound
			}

			poolDAO.Put(daoUser)
			return nil
		})
	})
	if err != nil {
		return errpkg.New("UPDATE_USER_TX_ERROR", "transaction error while updating user", err)
	}
	return nil
}

func (repo *User) List(
	ctx context.Context,
	pageSize uint,
	pageToken string,
) ([]*models.User, string, error) {
	const op = "repository.User.List"

	sqlStr := goqu.Select(
		"id", "email", "username", "password_hash", "created_at", "updated_at",
	).From("users")

	paginator := repo.paginator.WithDataset(sqlStr).
		WithColumns("created_at", "id").
		WithLimit(pageSize)
	itemsDAO, nextToken, err := paginator.Paginate(ctx, pageToken)
	if err != nil {
		return nil, "", errpkg.New("LIST_USERS_ERROR", "failed to paginate users", err)
	}

	// Берём срез из пула
	ptr := poolSliceUsers.Get().(*[]*models.User)
	users := *ptr
	// Очищаем содержимое, но оставляем capacity
	users = users[:0]

	for i := range itemsDAO {
		daoUser := &itemsDAO[i]
		modelUser := poolModel.Get().(*models.User)
		*modelUser = models.User{
			ID:        uint64(daoUser.ID.Int64),
			Email:     daoUser.Email,
			Username:  daoUser.Username,
			Password:  daoUser.Password,
			CreatedAt: daoUser.CreatedAt.Time,
			UpdatedAt: daoUser.UpdatedAt.Time,
		}
		users = append(users, modelUser)
	}

	// Скопируем users (к примеру, чтобы не держать pointer на слайс из пула)
	// Затем вернём слайс в пул
	result := make([]*models.User, len(users))
	copy(result, users)
	*ptr = users
	poolSliceUsers.Put(ptr)

	return result, nextToken, nil
}

func (repo *User) Delete(
	ctx context.Context,
	id string,
) error {
	const op = "repository.User.Delete"

	err := repo.Client.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		return repo.WithTx(ctx, func(tx pgx.Tx) error {
			stmt, args, err := repo.Dialect.Delete("users").
				Where(goqu.C("id").Eq(id)).
				Prepared(true).
				ToSQL()
			if err != nil {
				return errpkg.New("DELETE_USER_ERROR", "failed to build SQL query", err)
			}

			cmd, err := tx.Exec(ctx, stmt, args...)
			if err != nil {
				return errpkg.New("DELETE_USER_ERROR", "failed to delete user", err)
			}

			if cmd.RowsAffected() == 0 {
				return errpkg.ErrUserNotFound
			}

			return nil
		})
	})
	if err != nil {
		return errpkg.New("DELETE_USER_TX_ERROR", "transaction error while deleting user", err)
	}
	return nil
}
