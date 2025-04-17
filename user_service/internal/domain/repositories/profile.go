package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	dto "github.com/DENFNC/Zappy/user_service/internal/dto/profile"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *models.Profile) (uint32, error)
	GetByID(ctx context.Context, id uint32) (*models.Profile, error)
	List(ctx context.Context, params *dto.ListParams) (*dto.ListResult, error)
	Update(ctx context.Context, profile *models.Profile) (uint32, error)
	Delete(ctx context.Context, id uint32) (uint32, error)
}
