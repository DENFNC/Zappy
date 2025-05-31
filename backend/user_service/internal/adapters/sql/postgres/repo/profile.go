package repo

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type ProfileRepo struct {
	*postgres.Storage
}

func NewProfileRepo(
	db *postgres.Storage,
) *ProfileRepo {
	return &ProfileRepo{
		Storage: db,
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
		return "", err
	}

	var profileID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&profileID); err != nil {
		return "", err
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
		return nil, err
	}

	var profile models.Profile
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&profile.ProfileID,
		&profile.FirstName,
		&profile.LastName,
	); err != nil {
		return nil, err
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
		return nil, err
	}

	var profile models.Profile
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&profile.ProfileID,
		&profile.FirstName,
		&profile.LastName,
	); err != nil {
		return nil, err
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
		return "", err
	}

	var profileID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&profileID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", err
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
		return "", err
	}

	var profileID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&profileID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", err
	}

	return profileID, nil
}

func (r *ProfileRepo) List(ctx context.Context, params []any) ([]*models.Profile, error) {
	stmt, args, err := r.Dialect.Select(
		"profile_id",
		"user_id",
		"first_name",
		"last_name",
		"email",
		"phone",
	).From("profiles").
		Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*models.Profile
	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(
			&profile.ProfileID,
			&profile.FirstName,
			&profile.LastName,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, &profile)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}
