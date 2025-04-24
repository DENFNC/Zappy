package repo

import (
	"context"
	"fmt"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	psql "github.com/DENFNC/Zappy/user_service/internal/storage/postgres"
	"github.com/doug-martin/goqu/v9"
)

type ProfileRepo struct {
	*psql.Storage
	goqu *goqu.DialectWrapper
}

func NewProfileRepo(
	db *psql.Storage,
	goqu *goqu.DialectWrapper,
) *ProfileRepo {
	return &ProfileRepo{
		db,
		goqu,
	}
}
func (r *ProfileRepo) Create(ctx context.Context, profile *models.Profile) (string, error) {
	id := r.NewV7().String()

	stmt, args, err := r.goqu.
		Insert("profile").
		Returning(goqu.C("profile_id")).
		Rows(goqu.Record{
			"profile_id":   id,
			"auth_user_id": profile.AuthUserID,
			"first_name":   profile.FirstName,
			"last_name":    profile.LastName,
		}).
		Prepared(true).
		ToSQL()

	if err != nil {
		return "", fmt.Errorf("failed to build SQL: %w", err)
	}

	var profileID string
	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&profileID)
	if err != nil {
		return "", fmt.Errorf("failed to insert profile: %w", err)
	}

	return profileID, nil
}

func (r *ProfileRepo) GetByID(ctx context.Context, id string) (*models.Profile, error) {
	var profile models.Profile

	stmt, args, err := r.goqu.
		Select(
			"profile_id",
			"auth_user_id",
			"first_name",
			"last_name",
			"created_at",
			"updated_at",
		).
		From("profile").
		Where(goqu.C("profile_id").Eq(id)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	err = r.DB.QueryRow(ctx, stmt, args...).Scan(
		&profile.ProfileID,
		&profile.AuthUserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepo) List(
	ctx context.Context,
	item []any,
) (
	[]any,
	error,
) {
	panic("implement me!")
}

func (r *ProfileRepo) Update(ctx context.Context, profile *models.Profile) (string, error) {
	stmt, args, err := r.goqu.
		Update("profile").
		Returning(goqu.C("profile_id")).
		Set(goqu.Record{
			"first_name": profile.FirstName,
			"last_name":  profile.LastName,
			"updated_at": profile.UpdatedAt,
		}).
		Where(goqu.C("profile_id").Eq(profile.ProfileID)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return "", fmt.Errorf("failed to build SQL: %w", err)
	}

	var profileID string
	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&profileID)
	if err != nil {
		return "", fmt.Errorf("failed to update profile: %w", err)
	}

	return profileID, nil
}

func (r *ProfileRepo) Delete(ctx context.Context, id string) (string, error) {
	stmt, args, err := r.goqu.
		Delete("profile").
		Returning(goqu.C("profile_id")).
		Where(goqu.C("profile_id").Eq(id)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return "", fmt.Errorf("failed to build SQL: %w", err)
	}

	var profileID string
	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&profileID)
	if err != nil {
		return "", fmt.Errorf("failed to delete profile: %w", err)
	}

	return profileID, nil
}
