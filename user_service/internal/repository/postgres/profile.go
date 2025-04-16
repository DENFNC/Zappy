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

	err := r.DB.QueryRow(
		ctx,
		`SELECT profile_id,
				auth_user_id,
				first_name,
				last_name,
				created_at,
				updated_at
			FROM profile
		WHERE profile_id = $1`,
		id,
	).Scan(
		&profile.ProfileID,
		&profile.AuthUserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfileRepo) Update(ctx context.Context, profile *models.Profile) (uint32, error) {
	panic("implement me!")

}

func (r *ProfileRepo) Delete(ctx context.Context, id uint32) (uint32, error) {
	panic("implement me!")
}
