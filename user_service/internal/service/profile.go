package service

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
	emptyStringValue = ""
)

type Profile struct {
	log  *slog.Logger
	repo repositories.ProfileRepository
}

func NewProfile(
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
	authUserID string,
	firstName string,
	lastName string,
) (string, error) {
	const op = "service.Profile.Create"

	log := p.log.With("op", op)

	profile := models.NewProfile(
		authUserID,
		firstName,
		lastName,
	)

	profileID, err := p.repo.Create(ctx, profile)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return emptyStringValue, errpkg.New("CREATE_ERROR", "Couldn't create profile", err)
	}

	return profileID, nil
}

func (p *Profile) Delete(ctx context.Context, profileID string) (string, error) {
	const op = "service.Profile.Delete"

	log := p.log.With("op", op)

	profileID, err := p.repo.Delete(ctx, profileID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Profile not found")
			return emptyStringValue, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)

		return emptyStringValue, errpkg.New("DELETE_ERROR", "couldn't delete value", err)
	}

	return profileID, nil
}

func (p *Profile) GetByID(ctx context.Context, profileID string) (*models.Profile, error) {
	const op = "service.Profile.GetByID"

	log := p.log.With("op", op)

	profile, err := p.repo.GetByID(ctx, profileID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Profile not found")
			return nil, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)

		return nil, errpkg.New("GET_BY_ID_ERROR", "couldn't get value", err)
	}

	return profile, nil
}

func (p *Profile) List(context.Context, []any) ([]any, string, error) {
	panic("implement me")
}

func (p *Profile) Update(ctx context.Context, profileID string, firstName, lastName string) (string, error) {
	const op = "service.Profile.Update"

	log := p.log.With("op", op)

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
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Profile not found")
			return emptyStringValue, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)

		return emptyStringValue, errpkg.New("UPDATE_ERROR", "couldn't update value", err)
	}

	return profileID, nil
}
