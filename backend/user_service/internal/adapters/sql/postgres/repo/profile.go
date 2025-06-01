package repo

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/paginate"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type ProfileRepo struct {
	*postgres.Storage
	*paginate.Paginator[dao.ProfileDAO]
}

func NewProfileRepo(
	db *postgres.Storage,
	coder paginate.TokenCoder,
) *ProfileRepo {
	paginate, err := paginate.NewPaginator[dao.ProfileDAO](
		db.Client, db.Dialect, coder,
	)
	if err != nil {
		panic(err)
	}

	return &ProfileRepo{
		Storage:   db,
		Paginator: paginate,
	}
}

func (r *ProfileRepo) Create(ctx context.Context, profile *models.Profile) (string, error) {
	uid, _ := uuid.NewV7()
	stmt, args, err := r.Dialect.Insert("profiles").
		Rows(goqu.Record{
			"profile_id": uid.String(),
			"first_name": profile.FirstName,
			"last_name":  profile.LastName,
		}).
		Returning("profile_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", errpkg.New("PROFILE_CREATE_SQL_BUILD_ERROR", "failed to build create profile sql query", err)
	}

	var profileID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&profileID); err != nil {
		return "", errpkg.New("PROFILE_CREATE_QUERY_ERROR", "failed to execute create profile query", err)
	}

	return profileID, nil
}

func (r *ProfileRepo) GetByID(ctx context.Context, id string) (*models.Profile, error) {
	stmt, args, err := r.Dialect.Select(
		"profile_id",
		"user_id",
		"first_name",
		"last_name",
		"email",
		"phone",
	).From("profiles").
		Where(goqu.Ex{
			"profile_id": id,
		}).Prepared(true).ToSQL()
	if err != nil {
		return nil, errpkg.New("PROFILE_GETBYID_SQL_BUILD_ERROR", "failed to build get profile by id sql query", err)
	}

	var profile models.Profile
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&profile.ProfileID,
		&profile.FirstName,
		&profile.LastName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrNotFound
		}
		return nil, errpkg.New("PROFILE_GETBYID_QUERY_ERROR", "failed to execute get profile by id query", err)
	}

	return &profile, nil
}

func (r *ProfileRepo) GetByUserID(ctx context.Context, userID string) (*models.Profile, error) {
	stmt, args, err := r.Dialect.Select(
		"profile_id",
		"user_id",
		"first_name",
		"last_name",
		"email",
		"phone",
	).From("profiles").
		Where(goqu.Ex{
			"user_id": userID,
		}).Prepared(true).ToSQL()
	if err != nil {
		return nil, errpkg.New("PROFILE_GETBYUSERID_SQL_BUILD_ERROR", "failed to build get profile by user id sql query", err)
	}

	var profile models.Profile
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&profile.ProfileID,
		&profile.FirstName,
		&profile.LastName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrNotFound
		}
		return nil, errpkg.New("PROFILE_GETBYUSERID_QUERY_ERROR", "failed to execute get profile by user id query", err)
	}

	return &profile, nil
}

func (r *ProfileRepo) Update(ctx context.Context, profile *models.Profile) (string, error) {
	stmt, args, err := r.Dialect.Update("profiles").
		Set(goqu.Record{
			"first_name": profile.FirstName,
			"last_name":  profile.LastName,
		}).
		Where(goqu.Ex{
			"profile_id": profile.ProfileID,
		}).
		Returning("profile_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", errpkg.New("PROFILE_UPDATE_SQL_BUILD_ERROR", "failed to build update profile sql query", err)
	}

	var profileID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&profileID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", errpkg.New("PROFILE_UPDATE_QUERY_ERROR", "failed to execute update profile query", err)
	}

	return profileID, nil
}

func (r *ProfileRepo) Delete(ctx context.Context, id string) (string, error) {
	stmt, args, err := r.Dialect.Delete("profiles").
		Where(goqu.Ex{
			"profile_id": id,
		}).
		Returning("profile_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", errpkg.New("PROFILE_DELETE_SQL_BUILD_ERROR", "failed to build delete profile sql query", err)
	}

	var profileID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&profileID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", errpkg.New("PROFILE_DELETE_QUERY_ERROR", "failed to execute delete profile query", err)
	}

	return profileID, nil
}

func (repo *ProfileRepo) List(
	ctx context.Context,
	pageSize uint,
	pageToken string,
) ([]*models.Profile, string, error) {
	ds := repo.Dialect.Select(
		"profile_id",
		"user_id",
		"first_name",
		"last_name",
		"email",
		"phone",
	).From("profiles")

	repo.Paginator.WithDataset(ds).WithColumns("profile_id").WithLimit(pageSize)

	itemsDAO, nextPageToken, err := repo.Paginator.Paginate(
		ctx,
		pageToken,
	)
	if err != nil {
		return nil, "", errpkg.New("PROFILE_LIST_PAGINATE_ERROR", "failed to paginate profiles", err)
	}

	items := make([]*models.Profile, len(itemsDAO))
	for i, itemDAO := range itemsDAO {
		items[i] = &models.Profile{
			ProfileID:  itemDAO.ProfileID.String(),
			AuthUserID: itemDAO.AuthUserID.String(),
			FirstName:  itemDAO.FirstName.String,
			LastName:   itemDAO.LastName.String,
			CreatedAt:  itemDAO.CreatedAt.Time,
			UpdatedAt:  itemDAO.UpdatedAt.Time,
		}
	}

	return items, nextPageToken, nil
}
