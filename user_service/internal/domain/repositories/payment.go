package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
)

type PaymentRepository interface {
	Create(ctx context.Context, method *models.Payment) (string, error)
	GetByID(ctx context.Context, id string) (*models.Payment, error)
	GetByProfileID(ctx context.Context, profileID string) ([]models.Payment, error)
	SetDefault(ctx context.Context, methodID string, profileID string) (string, error)
	Update(ctx context.Context, method *models.Payment) (string, error)
	Delete(ctx context.Context, id string) (string, error)
}
