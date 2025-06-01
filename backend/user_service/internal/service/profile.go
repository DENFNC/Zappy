package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
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
		pageSize uint,
		pageToken string,
	) ([]*models.Profile, string, error)
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
	*slog.Logger
	ProfileRepository
}

func NewProfile(
	log *slog.Logger,
	repo ProfileRepository,
) *Profile {
	return &Profile{
		Logger:            log,
		ProfileRepository: repo,
	}
}

func (svc *Profile) CreateProfile(
	ctx context.Context,
	profile *models.Profile,
) (string, error) {
	const op = "service.Profile.Create"

	log := svc.Logger.With("op", op)

	profileID, err := svc.ProfileRepository.Create(ctx, profile)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return emptyStringValue, errpkg.New("CREATE_ERROR", "Couldn't create profile", err)
	}

	return profileID, nil
}

func (svc *Profile) DeleteProfile(
	ctx context.Context,
	profileID string,
) (string, error) {
	const op = "service.Profile.Delete"

	log := svc.Logger.With("op", op)

	profileID, err := svc.ProfileRepository.Delete(ctx, profileID)

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

func (svc *Profile) ProfileGetByID(
	ctx context.Context,
	profileID string,
) (*models.Profile, error) {
	const op = "service.Profile.GetByID"

	log := svc.Logger.With("op", op)

	profile, err := svc.ProfileRepository.GetByID(ctx, profileID)

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

func (svc *Profile) ListProfiles(
	ctx context.Context,
	pageSize uint,
	pageToken string,
) ([]*models.Profile, string, error) {
	const op = "service.Profile.List"

	log := svc.Logger.With("op", op)

	items, pageToken, err := svc.ProfileRepository.List(ctx, pageSize, pageToken)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, "", errpkg.New("LIST_ERROR", "couldn't list profiles", err)
	}

	return items, pageToken, nil
}

func (svc *Profile) UpdateProfile(
	ctx context.Context,
	profileID string,
	profile *models.Profile,
) (string, error) {
	const op = "service.Profile.Update"

	log := svc.Logger.With("op", op)

	profileID, err := svc.ProfileRepository.Update(
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
