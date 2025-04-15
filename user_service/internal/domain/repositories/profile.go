package repositories

import "github.com/DENFNC/Zappy/user_service/internal/domain/models"

type ProfileRepository interface {
	Create(profile *models.Profile) (int, error)
	GetByID(id int) (*models.Profile, error)
	GetByAuthUserID(authUserID int) (*models.Profile, error)
	Update(profile *models.Profile) error
	Delete(id int) error
}
