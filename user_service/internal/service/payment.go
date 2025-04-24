package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/domain/repositories"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type UUID = uuid.UUID

type Payment struct {
	log  *slog.Logger
	repo repositories.PaymentRepository
}

func NewPayment(
	log *slog.Logger,
	repo repositories.PaymentRepository,
) *Payment {
	return &Payment{
		log:  log,
		repo: repo,
	}
}

func (p *Payment) Create(
	ctx context.Context,
	profileID uint32,
	paymentToken string,
) (UUID, error) {
	const op = "service.Payment.Create"

	log := p.log.With("op", op)

	payment := models.NewPayment(profileID, paymentToken)

	payID, err := p.repo.Create(ctx, payment)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return uuid.Nil, err
	}

	return payID, nil
}

func (p *Payment) GetByID(
	ctx context.Context,
	paymentID uint32,
) (*models.Payment, error) {
	const op = "service.Payment.GetByID"

	log := p.log.With("op", op)

	payment, err := p.repo.GetByID(ctx, paymentID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return payment, nil
}

func (p *Payment) Update(
	ctx context.Context,
	paymentID uint32,
	profileID uint32,
	paymentToken string,
) (uint32, error) {
	const op = "service.Payment.Update"

	log := p.log.With("op", op)

	payment := models.Payment{
		PaymentID:    paymentID,
		ProfileID:    profileID,
		PaymentToken: paymentToken,
	}

	payID, err := p.repo.Update(ctx, &payment)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error(
				"Not found",
				slog.String("error", err.Error()),
			)
			return 0, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return 0, err
	}

	return payID, nil
}

func (p *Payment) Delete(
	ctx context.Context,
	paymentID uint32,
) (uint32, error) {
	const op = "service.Payment.Delete"

	log := p.log.With("op", op)

	payID, err := p.repo.Delete(ctx, paymentID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return 0, err
	}

	return payID, nil
}

func (p *Payment) List(
	ctx context.Context,
	profileID uint32,
) ([]models.Payment, error) {
	const op = "service.Payment.List"

	log := p.log.With("op", op)

	payments, err := p.repo.GetByProfileID(ctx, profileID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return payments, nil
}

func (p *Payment) SetDefault(
	ctx context.Context,
	paymentID uint32,
	profileID uint32,
) (uint32, error) {
	const op = "service.Payment.SetDefault"

	log := p.log.With("op", op)

	payID, err := p.repo.SetDefault(ctx, paymentID, profileID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return emptyValue, nil
	}

	return payID, nil
}
