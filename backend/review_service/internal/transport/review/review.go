package review

import (
	"context"

	"github.com/DENFNC/Zappy/review_service/internal/domain/models"
	"github.com/DENFNC/Zappy/review_service/proto/gen/go/common/v1"
	v1 "github.com/DENFNC/Zappy/review_service/proto/gen/go/review/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Review interface {
	CreateReview(
		ctx context.Context,
		input *models.Review,
	) (string, error)
	// GetReviewByID(
	// 	ctx context.Context,
	// 	id string,
	// ) (Review, error)
	// ListReviews(ctx context.Context, filter ReviewFilter) ([]Review, error)
	// UpdateReview(ctx context.Context, id string, input UpdateReviewInput) (Review, error)
	// DeleteReview(
	// 	ctx context.Context,
	// 	id string,
	// ) error
}

type serverAPI struct {
	v1.UnimplementedReviewServiceServer
	Review
}

func New(svc Review) *serverAPI {
	return &serverAPI{
		Review: svc,
	}
}

func (api *serverAPI) GRPCRegister(grpc *grpc.Server) {
	v1.RegisterReviewServiceServer(grpc, api)
}

func (api serverAPI) CreateReview(
	ctx context.Context,
	req *v1.CreateReviewRequest,
) (*v1.CreateReviewResponse, error) {
	uid, err := api.Review.CreateReview(ctx,
		&models.Review{
			ProductID: req.GetProductId(),
			ProfileID: req.GetProfileId(),
			Comment:   req.GetComment(),
			Rating:    req.GetRating(),
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &v1.CreateReviewResponse{
		ReviewId: &common.ResourceID{
			Id: uid,
		},
	}, nil
}

func (api *serverAPI) GetReview(
	ctx context.Context,
	req *v1.GetReviewRequest,
) (*v1.GetReviewResponse, error) {
	panic("implement me!")
}

func (api *serverAPI) ListReviews(
	ctx context.Context,
	req *v1.ListReviewsRequest,
) (*v1.ListReviewsResponse, error) {
	panic("implement me!")
}

func (api *serverAPI) DeleteReview(
	ctx context.Context,
	req *v1.DeleteReviewRequest,
) (*v1.DeleteReviewResponse, error) {
	panic("implement me!")
}
