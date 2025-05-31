package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/jackc/pgx/v5"
)

type PaymentRepository interface {
	Create(
		ctx context.Context,
		method *models.Payment,
	) (string, error)
	GetByID(
		ctx context.Context,
		uid string,
	) (*models.Payment, error)
	GetByProfileID(
		ctx context.Context,
		profileID string,
	) ([]models.Payment, error)
	SetDefault(
		ctx context.Context,
		methodID string,
		profileID string,
	) error
	Update(
		ctx context.Context,
		uid string,
		payment *models.Payment,
	) (string, error)
	Delete(
		ctx context.Context,
		uid string,
	) (string, error)
}

type Payment struct {
	log  *slog.Logger
	repo PaymentRepository
}

func NewPayment(
	log *slog.Logger,
	repo PaymentRepository,
) *Payment {
	return &Payment{
		log:  log,
		repo: repo,
	}
}

func (p *Payment) CreatePayment(
	ctx context.Context,
	payment *models.Payment,
) (string, error) {
	const op = "service.Payment.Create"

	log := p.log.With("op", op)

	payID, err := p.repo.Create(ctx, payment)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return payID, nil
}

func (p *Payment) PaymentGetByID(
	ctx context.Context,
	paymentID string,
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

func (p *Payment) UpdatePayment(
	ctx context.Context,
	uid string,
	payment *models.Payment,
) (string, error) {
	const op = "service.Payment.Update"

	log := p.log.With("op", op)

	payID, err := p.repo.Update(ctx, uid, payment)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error(
				"Not found",
				slog.String("error", err.Error()),
			)
			return "", errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return payID, nil
}

func (p *Payment) DeletePayment(
	ctx context.Context,
	paymentID string,
) (string, error) {
	const op = "service.Payment.Delete"

	log := p.log.With("op", op)

	payID, err := p.repo.Delete(ctx, paymentID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return payID, nil
}

func (p *Payment) ListPayments(
	ctx context.Context,
	profileID string,
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

func (p *Payment) SetDefaultPayment(
	ctx context.Context,
	payment *models.Payment,
) error {
	const op = "service.Payment.SetDefault"

	log := p.log.With("op", op)

	err := p.repo.SetDefault(ctx, payment.PaymentID, payment.ProfileID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil
	}

	return nil
}
