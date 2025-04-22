package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
)

type PaymentRepository interface {
	Create(ctx context.Context, method *models.Payment) (uint32, error)
	GetByID(ctx context.Context, id uint32) (*models.Payment, error)
	GetByProfileID(ctx context.Context, profileID uint32) ([]models.Payment, error)
	SetDefault(ctx context.Context, methodID uint32, profileID uint32) (uint32, error)
	Update(ctx context.Context, method *models.Payment) (uint32, error)
	Delete(ctx context.Context, id uint32) (uint32, error)
}
