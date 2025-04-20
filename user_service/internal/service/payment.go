package service

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/domain/repositories"
)

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
	isDefault bool,
) (uint32, error) {
	panic("implement me!")
}

func (p *Payment) GetByID(
	ctx context.Context,
	paymentID uint32,
) (*models.Payment, error) {
	panic("implement me!")
}

func (p *Payment) Update(
	ctx context.Context,
	paymentID uint32,
	profileID uint32,
	paymentToken string,
	isDefault bool,
) (uint32, error) {
	panic("implement me!")
}

func (p *Payment) Delete(
	ctx context.Context,
	paymentID uint32,
) (uint32, error) {
	panic("implement me!")
}

func (p *Payment) List(
	ctx context.Context,
	profileID uint32,
) ([]*models.Payment, error) {
	panic("implement me!")
}

func (p *Payment) SetDefault(
	ctx context.Context,
	paymentID uint32,
	profileID uint32,
) (uint32, error) {
	panic("implement me!")
}
