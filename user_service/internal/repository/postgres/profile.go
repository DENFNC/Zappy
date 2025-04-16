package psqlrepoprofile

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

func New(
	db *psql.Storage,
	goqu *goqu.DialectWrapper,
) *ProfileRepo {
	return &ProfileRepo{
		db,
		goqu,
	}
}

func (r *ProfileRepo) Create(ctx context.Context, profile *models.Profile) (uint32, error) {
	var profileID uint32

	stmt, args, err := r.goqu.Insert("profile").
		Returning(goqu.C("profile_id")).
		Rows(goqu.Record{
			"auth_user_id": profile.AuthUserID,
			"first_name":   profile.FirstName,
			"last_name":    profile.LastName,
		}).
		Prepared(true).
		ToSQL()

	fmt.Println(stmt)

	if err != nil {
		return 0, fmt.Errorf("failed to build SQL: %w", err)
	}

	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&profileID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert profile: %w", err)
	}

	return profileID, nil
}

func (r *ProfileRepo) GetByID(
	ctx context.Context,
	id uint32,
) (*models.Profile, error) {
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
		return nil, err
	}

	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(
		&profile.ProfileID,
		&profile.AuthUserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfileRepo) Update(ctx context.Context, profile *models.Profile) (uint32, error) {
	stmt, args, err := r.goqu.Update("profile").Returning(goqu.C("profile_id")).Set(
		goqu.Record{
			"first_name": profile.FirstName,
			"last_name":  profile.LastName,
			"updated_at": profile.UpdatedAt,
		},
	).Where(goqu.C("profile_id").Eq(profile.ProfileID)).Prepared(true).ToSQL()

	if err != nil {
		return 0, err
	}

	var profileID uint32
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(&profileID); err != nil {
		return 0, err
	}

	return profileID, nil
}

func (r *ProfileRepo) Delete(ctx context.Context, id uint32) (uint32, error) {
	panic("implement me!")
}
