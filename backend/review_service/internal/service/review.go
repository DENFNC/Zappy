package service

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/review_service/internal/domain/models"
)

type ReviewRepo interface {
	Create(
		ctx context.Context,
		review *models.Review,
	) (string, error)
	GetByID(
		ctx context.Context,
		uid string,
	) (*models.Review, error)
}

type Review struct {
	*slog.Logger
	ReviewRepo
}

func New(log *slog.Logger, repo ReviewRepo) *Review {
	return &Review{
		Logger:     log,
		ReviewRepo: repo,
	}
}

func (svc *Review) CreateReview(
	ctx context.Context,
	input *models.Review,
) (string, error) {
	const op = "service.Review.CreateReview"

	log := svc.Logger.With("op", op)

	uid, err := svc.ReviewRepo.Create(
		ctx,
		&models.Review{
			ProductID: input.ProductID,
			ProfileID: input.ProfileID,
			Rating:    input.Rating,
			Comment:   input.Comment,
		},
	)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return uid, nil
}

func (svc *Review) GetReviewByID(
	ctx context.Context,
	uid string,
) (*models.Review, error) {
	const op = "service.Review.GetReviewByID"

	log := svc.Logger.With("op", op)

	data, err := svc.ReviewRepo.GetByID(ctx, uid)

	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return data, nil
}
