package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *models.Profile) (string, error)
	GetByID(ctx context.Context, id string) (*models.Profile, error)
	List(ctx context.Context, params []any) ([]any, error)
	Update(ctx context.Context, profile *models.Profile) (string, error)
	Delete(ctx context.Context, id string) (string, error)
}
