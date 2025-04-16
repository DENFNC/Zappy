package profileservice

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/domain/repositories"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/jackc/pgx/v5"
)

const (
	emptyValue = 0
)

type Profile struct {
	log  *slog.Logger
	repo repositories.ProfileRepository
}

func New(
	log *slog.Logger,
	repo repositories.ProfileRepository,
) *Profile {
	return &Profile{
		log:  log,
		repo: repo,
	}
}

func (p *Profile) Create(
	ctx context.Context,
	authUserID uint32,
	firstName string,
	lastName string,
) (uint32, error) {
	profile := models.NewProfile(
		authUserID,
		firstName,
		lastName,
	)

	profileID, err := p.repo.Create(ctx, profile)
	if err != nil {
		return emptyValue, err
	}

	return profileID, nil
}

func (p *Profile) Delete(ctx context.Context, profileID uint32) (uint32, error) {
	const op = "service.Profile.Delete"

	log := p.log.With("op", op)

	profileID, err := p.repo.Delete(ctx, profileID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Not found")
			return emptyValue, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.Any("error", err),
		)

		return emptyValue, errpkg.New("DELETE_ERROR", "couldn't delete value", err)
	}

	return profileID, nil
}

func (p *Profile) GetByID(ctx context.Context, profileID uint32) (*models.Profile, error) {
	profile, err := p.repo.GetByID(ctx, profileID)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (p *Profile) List(ctx context.Context) ([]*models.Profile, error) {
	panic("Implement me!")
}

func (p *Profile) Update(ctx context.Context, profileID uint32, firstName, lastName string) (uint32, error) {
	profile := &models.Profile{
		ProfileID: profileID,
		FirstName: firstName,
		LastName:  lastName,
		UpdatedAt: time.Now(),
	}

	profileID, err := p.repo.Update(
		ctx,
		profile,
	)
	if err != nil {
		return emptyValue, err
	}

	return profileID, nil
}
