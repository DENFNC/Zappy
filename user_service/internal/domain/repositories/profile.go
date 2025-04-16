package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *models.Profile) (uint32, error)
	GetByID(ctx context.Context, id uint32) (*models.Profile, error)
	Update(ctx context.Context, profile *models.Profile) (uint32, error)
	Delete(ctx context.Context, id uint32) (uint32, error)
}
