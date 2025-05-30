package service

import (
	"context"

	"github.com/DENFNC/Zappy/review_service/internal/domain/models"
)

type ReviewRepo interface {
	Create(
		ctx context.Context,
		review *models.Review,
	) (string, error)
}

type Review struct {
	ReviewRepo
}

func New(repo ReviewRepo) *Review {
	return &Review{
		ReviewRepo: repo,
	}
}

func (svc *Review) CreateReview(
	ctx context.Context,
	input *models.Review,
) (string, error) {
	uid, err := svc.ReviewRepo.Create(
		ctx,
		&models.Review{
			ProductID: input.ProductID,
			ProfileID: input.ProfileID,
			Rating:    input.Rating,
		},
	)
	if err != nil {
		return "", err
	}

	return uid, nil
}
