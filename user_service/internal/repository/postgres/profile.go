package psqlrepoprofile

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	psql "github.com/DENFNC/Zappy/user_service/internal/storage/postgres"
)

type ProfileRepo struct {
	*psql.Storage
}

func New(
	db *psql.Storage,
) *ProfileRepo {
	return &ProfileRepo{
		db,
	}
}

func (r *ProfileRepo) Create(ctx context.Context, profile *models.Profile) (uint64, error) {
	var profileID uint64

	err := r.DB.QueryRow(
		ctx,
		`INSERT INTO profile (
			auth_user_id,
			first_name,
			last_name
		) VALUES ($1, $2, $3)
		RETURNING profile_id;
		`,
		profile.AuthUserID,
		profile.FirstName,
		profile.LastName,
	).Scan(&profileID)

	if err != nil {
		return 0, err
	}

	return profileID, nil
}

func (r *ProfileRepo) GetByID(
	ctx context.Context,
	id uint64,
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

func (r *ProfileRepo) Update(ctx context.Context, profile *models.Profile) (uint64, error) {
	panic("implement me!")
}

func (r *ProfileRepo) Delete(ctx context.Context, id int) (uint64, error) {
	panic("implement me!")
}
