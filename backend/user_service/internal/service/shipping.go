package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/domain/repositories"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const emptyStringAddr = ""

type ShippingService struct {
	log  *slog.Logger
	repo repositories.ShippingRepository
}

func NewShipping(log *slog.Logger, repo repositories.ShippingRepository) *ShippingService {
	return &ShippingService{
		log:  log,
		repo: repo,
	}
}

func (s *ShippingService) Create(ctx context.Context, address *models.Shipping) (string, error) {
	const op = "service.ShippingService.Create"

	log := s.log.With("op", op)

	addrID, err := s.repo.Create(ctx, address)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return emptyStringAddr, errpkg.ErrConstraint
			case "23505":
				return emptyStringAddr, errpkg.ErrUniqueViolation
			}
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return emptyStringAddr, err
	}

	return addrID, nil
}

func (s *ShippingService) GetByID(ctx context.Context, id string) (*models.Shipping, error) {
	const op = "service.ShippingService.GetByID"

	log := s.log.With("op", op)

	address, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.String("address_id", id),
		)
		return nil, err
	}

	return address, nil
}

func (s *ShippingService) ListByProfile(ctx context.Context, profileID string) ([]models.Shipping, error) {
	const op = "service.ShippingService.ListByProfile"

	log := s.log.With("op", op)

	addresses, err := s.repo.GetByProfileID(ctx, profileID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.String("profile_id", profileID),
		)
		return nil, err
	}

	return addresses, nil
}

func (s *ShippingService) Update(ctx context.Context, id string, address *models.Shipping) (string, error) {
	const op = "service.ShippingService.Update"

	log := s.log.With("op", op)

	addrID, err := s.repo.UpdateAddress(
		ctx,
		id,
		address,
	)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return emptyStringAddr, err
	}

	return addrID, nil
}

func (s *ShippingService) SetDefault(ctx context.Context, addressID, profileID string) (string, error) {
	const op = "service.ShippingService.SetDefault"

	log := s.log.With("op", op)

	if err := s.repo.SetDefault(ctx, addressID, profileID); err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.String("address_id", addressID),
			slog.String("profile_id", profileID),
		)
		return emptyStringAddr, err
	}

	return addressID, nil
}

func (s *ShippingService) Delete(ctx context.Context, id string) (string, error) {
	const op = "service.ShippingService.Delete"

	log := s.log.With("op", op)

	addrID, err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyStringAddr, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.String("address_id", id),
		)
		return emptyStringAddr, err
	}

	return addrID, nil
}
