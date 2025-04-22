package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
)

type ShippingRepository interface {
	Create(
		ctx context.Context,
		address *models.Shipping,
	) (uint32, error)
	GetByID(
		ctx context.Context,
		id uint32,
	) (*models.Shipping, error)
	GetByProfileID(
		ctx context.Context,
		profileID uint32,
	) ([]models.Shipping, error)
	UpdateAddress(
		ctx context.Context,
		id uint32,
		address *models.Shipping,
	) (uint32, error)
	SetDefault(
		ctx context.Context,
		addressID, profileID uint32,
	) error
	Delete(
		ctx context.Context,
		id uint32,
	) (uint32, error)
}
