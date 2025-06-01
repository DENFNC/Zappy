package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
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
	List(
		ctx context.Context,
		profileID string,
		pageSize uint,
		pageToken string,
	) ([]*models.Payment, string, error)
	Delete(
		ctx context.Context,
		uid string,
	) (string, error)
}

type Payment struct {
	*slog.Logger
	PaymentRepository
}

func NewPayment(
	log *slog.Logger,
	repo PaymentRepository,
) *Payment {
	return &Payment{
		Logger:            log,
		PaymentRepository: repo,
	}
}

func (svc *Payment) CreatePayment(
	ctx context.Context,
	payment *models.Payment,
) (string, error) {
	const op = "service.Payment.Create"

	log := svc.Logger.With("op", op)

	payID, err := svc.PaymentRepository.Create(ctx, payment)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", errpkg.New("PAYMENT_CREATE_ERROR", "Failed to create payment", err)
	}

	return payID, nil
}

func (svc *Payment) PaymentGetByID(
	ctx context.Context,
	paymentID string,
) (*models.Payment, error) {
	const op = "service.Payment.GetByID"

	log := svc.Logger.With("op", op)

	payment, err := svc.PaymentRepository.GetByID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Payment not found")
			return nil, errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, errpkg.New("PAYMENT_GET_ERROR", "Failed to get payment", err)
	}

	return payment, nil
}

func (svc *Payment) UpdatePayment(
	ctx context.Context,
	uid string,
	payment *models.Payment,
) (string, error) {
	const op = "service.Payment.Update"

	log := svc.Logger.With("op", op)

	payID, err := svc.PaymentRepository.Update(ctx, uid, payment)
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

func (svc *Payment) DeletePayment(
	ctx context.Context,
	paymentID string,
) (string, error) {
	const op = "service.Payment.Delete"

	log := svc.Logger.With("op", op)

	payID, err := svc.PaymentRepository.Delete(ctx, paymentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Payment not found")
			return "", errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", errpkg.New("PAYMENT_DELETE_ERROR", "Failed to delete payment", err)
	}

	return payID, nil
}

func (svc *Payment) ListPayments(
	ctx context.Context,
	profileID string,
	pageSize uint,
	pageToken string,
) ([]*models.Payment, string, error) {
	const op = "service.Payment.ListPayments"

	log := svc.Logger.With("op", op)

	items, pageToken, err := svc.PaymentRepository.List(
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
		return nil, "", errpkg.New("PAYMENT_LIST_ERROR", "Failed to list payments", err)
	}

	return items, pageToken, nil
}

func (svc *Payment) SetDefaultPayment(
	ctx context.Context,
	payment *models.Payment,
) error {
	const op = "service.Payment.SetDefault"

	log := svc.Logger.With("op", op)

	err := svc.PaymentRepository.SetDefault(ctx, payment.PaymentID, payment.ProfileID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Payment not found")
			return errpkg.ErrNotFound
		}
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return errpkg.New("PAYMENT_SETDEFAULT_ERROR", "Failed to set default payment", err)
	}

	return nil
}
