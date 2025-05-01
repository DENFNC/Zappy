package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
)

type ShippingRepository interface {
	Create(
		ctx context.Context,
		address *models.Shipping,
	) (string, error)
	GetByID(
		ctx context.Context,
		id string,
	) (*models.Shipping, error)
	GetByProfileID(
		ctx context.Context,
		profileID string,
	) ([]models.Shipping, error)
	UpdateAddress(
		ctx context.Context,
		id string,
		address *models.Shipping,
	) (string, error)
	SetDefault(
		ctx context.Context,
		addressID, profileID string,
	) error
	Delete(
		ctx context.Context,
		id string,
	) (string, error)
}
