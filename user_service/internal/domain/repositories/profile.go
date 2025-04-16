package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *models.Profile) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*models.Profile, error)
	Update(ctx context.Context, profile *models.Profile) (uint64, error)
	Delete(ctx context.Context, id int) (uint64, error)
}
