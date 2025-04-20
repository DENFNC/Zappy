package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	repo "github.com/DENFNC/Zappy/user_service/internal/repository/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ShippingService struct {
	repo *repo.ShippingRepo
	log  *slog.Logger
}

func NewShipping(repo *repo.ShippingRepo, log *slog.Logger) *ShippingService {
	return &ShippingService{
		repo: repo,
		log:  log,
	}
}

func (s *ShippingService) Create(ctx context.Context, address *models.Shipping) (uint32, error) {
	const op = "service.ShippingService.Create"

	log := s.log.With("op", op)

	addrID, err := s.repo.Create(ctx, address)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return emptyValue, errpkg.ErrConstraint
			case "23505":
				return emptyValue, errpkg.ErrUniqueViolation
			}
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return emptyValue, err
	}

	return addrID, nil
}

func (s *ShippingService) GetByID(ctx context.Context, id int) (*models.Shipping, error) {
	const op = "service.ShippingService.GetByID"

	log := s.log.With("op", op)

	address, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.Int("address_id", id),
		)
		return nil, err
	}

	return address, nil
}

func (s *ShippingService) ListByProfile(ctx context.Context, profileID uint32) ([]*models.Shipping, error) {
	const op = "service.ShippingService.ListByProfile"

	log := s.log.With("op", op)

	addresses, err := s.repo.GetByProfileID(ctx, int(profileID))
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.Uint64("profile_id", uint64(profileID)),
		)
		return nil, err
	}

	return addresses, nil
}

func (s *ShippingService) SetDefault(ctx context.Context, addressID, profileID uint32) (uint32, error) {
	const op = "service.ShippingService.SetDefault"

	log := s.log.With("op", op)

	if err := s.repo.SetDefault(ctx, addressID, profileID); err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.Uint64("address_id", uint64(addressID)),
			slog.Uint64("profile_id", uint64(profileID)),
		)
		return emptyValue, err
	}

	return uint32(addressID), nil
}

func (s *ShippingService) Delete(ctx context.Context, id uint32) (uint32, error) {
	const op = "service.ShippingService.Delete"

	log := s.log.With("op", op)

	addrID, err := s.repo.Delete(ctx, int(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyValue, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.Uint64("address_id", uint64(id)),
		)
		return emptyValue, err
	}

	return addrID, nil
}
