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
	List(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
	) ([]models.Review, string, error)
	Update(
		ctx context.Context,
		uid string,
		comment string,
	) error
	DeleteByID(
		ctx context.Context,
		uid string,
	) error
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

func (svc *Review) ListReviews(
	ctx context.Context,
	pageSize uint32,
	pageToken string,
) ([]models.Review, string, error) {
	const op = "service.Review.ListReviews"

	log := svc.Logger.With("op", op)

	items, pageToken, err := svc.ReviewRepo.List(ctx, pageSize, pageToken)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, "", err
	}

	return items, pageToken, nil
}

func (svc *Review) UpdateReview(
	ctx context.Context,
	uid string,
	comment string,
) error {
	const op = "service.Review.ListReviews"

	log := svc.Logger.With("op", op)

	if err := svc.ReviewRepo.Update(ctx, uid, comment); err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (svc *Review) DeleteReviewByID(
	ctx context.Context,
	uid string,
) error {
	const op = "service.Review.DeleteReviewByID"

	log := svc.Logger.With("op", op)

	if err := svc.ReviewRepo.DeleteByID(ctx, uid); err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}
