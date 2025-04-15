package repositories

import "github.com/DENFNC/Zappy/user_service/internal/domain/models"

type PaymentMethodRepository interface {
	Create(method *models.PaymentMethod) (int, error)
	GetByID(id int) (*models.PaymentMethod, error)
	GetByProfileID(profileID int) ([]*models.PaymentMethod, error)
	SetDefault(methodID int, profileID int) error
	Delete(id int) error
}
