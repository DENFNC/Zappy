package review

import (
	"context"

	v1 "github.com/DENFNC/Zappy/review_service/proto/gen/go/review/v1"
)

type serverAPI struct {
	v1.UnimplementedReviewServiceServer
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
