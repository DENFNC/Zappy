package dto

import "github.com/DENFNC/Zappy/user_service/internal/domain/models"

type ListResult struct {
	Items         []*models.Profile
	NextPageToken string
}
