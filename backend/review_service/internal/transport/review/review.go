package review

import (
	"context"

	"github.com/DENFNC/Zappy/review_service/internal/domain/models"
	"github.com/DENFNC/Zappy/review_service/proto/gen/go/common/v1"
	v1 "github.com/DENFNC/Zappy/review_service/proto/gen/go/review/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Review interface {
	CreateReview(
		ctx context.Context,
		input *models.Review,
	) (string, error)
	GetReviewByID(
		ctx context.Context,
		uid string,
	) (*models.Review, error)
	ListReviews(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
	) ([]models.Review, string, error)
	UpdateReview(
		ctx context.Context,
		uid string,
		comment string,
	) error
	DeleteReviewByID(
		ctx context.Context,
		uid string,
	) error
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

func (api *serverAPI) HTTPRegister(
	ctx context.Context,
	mux *runtime.ServeMux,
) {
	v1.RegisterReviewServiceHandlerServer(ctx, mux, api)
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
	data, err := api.Review.GetReviewByID(ctx,
		req.ReviewId.GetId(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &v1.GetReviewResponse{
		Review: &common.Review{
			ReviewId:  data.ReviewID,
			ProductId: data.ProductID,
			ProfileId: data.ProductID,
			Rating:    data.Rating,
			Comment:   data.Comment,
			CreatedAt: timestamppb.New(data.CreatedAt),
			UpdatedAt: timestamppb.New(data.UpdatedAt),
		},
	}, nil
}

func (api *serverAPI) ListReviews(
	ctx context.Context,
	req *v1.ListReviewsRequest,
) (*v1.ListReviewsResponse, error) {
	afterPage := true
	items, pageToken, err := api.Review.ListReviews(ctx,
		req.Pagination.GetPageSize(),
		req.Pagination.GetPageToken(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}
	if pageToken == "" {
		afterPage = false
	}

	reviews := make([]*common.Review, len(items))
	for i, item := range items {
		reviews[i] = &common.Review{
			ReviewId:  item.ReviewID,
			ProductId: item.ProductID,
			ProfileId: item.ProfileID,
			Rating:    item.Rating,
			Comment:   item.Comment,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		}
	}

	return &v1.ListReviewsResponse{
		Reviews: reviews,
		Pagination: &common.PaginationResponse{
			PageSize:  uint32(len(items)),
			PageToken: pageToken,
			AfterPage: afterPage,
		},
	}, nil
}

func (api *serverAPI) UpdateReview(
	ctx context.Context,
	req *v1.UpdateReviewRequest,
) (*v1.UpdateReviewResponse, error) {
	if err := api.Review.UpdateReview(ctx,
		req.ReviewId.GetId(),
		req.GetComment(),
	); err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &v1.UpdateReviewResponse{}, nil
}

func (api *serverAPI) DeleteReview(
	ctx context.Context,
	req *v1.DeleteReviewRequest,
) (*v1.DeleteReviewResponse, error) {
	if err := api.Review.DeleteReviewByID(ctx,
		req.ReviewId.GetId(),
	); err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &v1.DeleteReviewResponse{}, nil
}
