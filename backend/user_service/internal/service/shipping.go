package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const emptyStringAddr = ""

type ShippingRepository interface {
	Create(
		ctx context.Context,
		address *models.Shipping,
	) (string, error)
	GetByID(
		ctx context.Context,
		id string,
	) (*models.Shipping, error)
	GetByProfileID(
		ctx context.Context,
		profileID string,
	) ([]models.Shipping, error)
	UpdateAddress(
		ctx context.Context,
		id string,
		address *models.Shipping,
	) (string, error)
	SetDefault(
		ctx context.Context,
		addressID, profileID string,
	) error
	Delete(
		ctx context.Context,
		id string,
	) (string, error)
	List(
		ctx context.Context,
		profileID string,
		pageSize uint,
		pageToken string,
	) ([]*models.Shipping, string, error)
}

type Shipping struct {
	*slog.Logger
	ShippingRepository
}

func NewShipping(
	log *slog.Logger,
	repo ShippingRepository,
) *Shipping {
	return &Shipping{
		Logger:             log,
		ShippingRepository: repo,
	}
}

func (svc *Shipping) Create(
	ctx context.Context,
	address *models.Shipping,
) (string, error) {
	const op = "service.ShippingService.Create"

	log := svc.Logger.With("op", op)

	addrID, err := svc.ShippingRepository.Create(ctx, address)
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
		return emptyStringAddr, errpkg.New("SHIPPING_CREATE_ERROR", "Failed to create shipping address", err)
	}

	return addrID, nil
}

func (svc *Shipping) GetByID(
	ctx context.Context,
	id string,
) (*models.Shipping, error) {
	const op = "service.ShippingService.GetByID"

	log := svc.Logger.With("op", op)

	address, err := svc.ShippingRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Shipping address not found", slog.String("address_id", id))
			return nil, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.String("address_id", id),
		)
		return nil, errpkg.New("SHIPPING_GET_ERROR", "Failed to get shipping address", err)
	}

	return address, nil
}

func (svc *Shipping) ListByProfile(
	ctx context.Context,
	profileID string,
) ([]models.Shipping, error) {
	const op = "service.ShippingService.ListByProfile"

	log := svc.Logger.With("op", op)

	addresses, err := svc.ShippingRepository.GetByProfileID(ctx, profileID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.String("profile_id", profileID),
		)
		return nil, errpkg.New("SHIPPING_LIST_ERROR", "Failed to list shipping addresses", err)
	}

	return addresses, nil
}

func (svc *Shipping) Update(
	ctx context.Context,
	id string,
	address *models.Shipping,
) (string, error) {
	const op = "service.ShippingService.Update"

	log := svc.Logger.With("op", op)

	addrID, err := svc.ShippingRepository.UpdateAddress(
		ctx,
		id,
		address,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Shipping address not found", slog.String("address_id", id))
			return emptyStringAddr, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return emptyStringAddr, errpkg.New("SHIPPING_UPDATE_ERROR", "Failed to update shipping address", err)
	}

	return addrID, nil
}

func (svc *Shipping) SetDefault(
	ctx context.Context,
	addressID, profileID string,
) (string, error) {
	const op = "service.ShippingService.SetDefault"

	log := svc.Logger.With("op", op)

	if err := svc.ShippingRepository.SetDefault(ctx, addressID, profileID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Shipping address not found",
				slog.String("address_id", addressID),
				slog.String("profile_id", profileID))
			return emptyStringAddr, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
			slog.String("address_id", addressID),
			slog.String("profile_id", profileID),
		)
		return emptyStringAddr, errpkg.New("SHIPPING_SETDEFAULT_ERROR", "Failed to set default shipping address", err)
	}

	return addressID, nil
}

func (svc *Shipping) ListShipping(
	ctx context.Context,
	profileID string,
	pageSize uint,
	pageToken string,
) ([]*models.Shipping, string, error) {
	const op = "service.Shipping.ListShipping"

	log := svc.Logger.With("op", op)

	items, pageToken, err := svc.ShippingRepository.List(
		ctx,
		profileID,
		pageSize,
		pageToken,
	)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, "", errpkg.New("SHIPPING_LIST_ERROR", "Failed to list shipping addresses", err)
	}

	return items, pageToken, nil
}
