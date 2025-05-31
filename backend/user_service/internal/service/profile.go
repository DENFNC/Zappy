package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/jackc/pgx/v5"
)

const (
	emptyStringValue = ""
)

type ProfileRepository interface {
	Create(
		ctx context.Context,
		profile *models.Profile,
	) (string, error)
	GetByID(
		ctx context.Context,
		uid string,
	) (*models.Profile, error)
	List(
		ctx context.Context,
		params []any,
	) ([]*models.Profile, error)
	Update(
		ctx context.Context,
		profile *models.Profile,
	) (string, error)
	Delete(
		ctx context.Context,
		uid string,
	) (string, error)
}

type Profile struct {
	log  *slog.Logger
	repo ProfileRepository
}

func NewProfile(
	log *slog.Logger,
	repo ProfileRepository,
) *Profile {
	return &Profile{
		log:  log,
		repo: repo,
	}
}

func (p *Profile) Create(
	ctx context.Context,
	profile *models.Profile,
) (string, error) {
	const op = "service.Profile.Create"

	log := p.log.With("op", op)

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

func (p *Profile) List(ctx context.Context, params []any) ([]any, string, error) {
	const op = "service.Profile.List"

	log := p.log.With("op", op)

	profiles, err := p.repo.List(ctx, params)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, "", errpkg.New("LIST_ERROR", "couldn't list profiles", err)
	}

	var anyProfiles []any
	for _, profile := range profiles {
		anyProfiles = append(anyProfiles, profile)
	}

	return anyProfiles, "", nil
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
