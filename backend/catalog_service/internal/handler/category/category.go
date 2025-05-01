package category

import (
	"context"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/catalog_service/internal/errors"
	v1 "github.com/DENFNC/Zappy/catalog_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Category interface {
	Create(
		ctx context.Context,
		name, parentID string,
	) (string, error)
	List(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
	) ([]models.Category, string, error)
	Delete(
		ctx context.Context,
		categoryID string,
	) error
}

type serverAPI struct {
	v1.UnimplementedCategoryServiceServer
	svc Category
}

func New(serv Category) *serverAPI {
	return &serverAPI{
		svc: serv,
	}
}

func (api *serverAPI) Register(grpc *grpc.Server) {
	v1.RegisterCategoryServiceServer(grpc, api)
}

func (api *serverAPI) CreateCategory(
	ctx context.Context,
	req *v1.CreateCategoryRequest,
) (*v1.CreateCategoryResponse, error) {
	categoryID, err := api.svc.Create(ctx, req.GetName(), req.GetParentId())
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.CreateCategoryResponse{
		CategoryId: &v1.ResourceID{
			Id: categoryID,
		},
	}, nil
}

func (api *serverAPI) ListCategories(
	ctx context.Context,
	req *v1.ListCategoriesRequest,
) (*v1.ListCategoriesResponse, error) {
	afterPage := true
	items, nextPageToken, err := api.svc.List(
		ctx,
		req.Pagination.GetPageSize(),
		req.Pagination.GetPageToken(),
	)
	if nextPageToken == "" {
		afterPage = false
	}
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	v1Category := make([]*v1.Category, len(items))
	for i, item := range items {
		v1Category[i] = &v1.Category{
			Id:        item.CategoryID,
			Name:      item.CategoryName,
			ParentId:  *item.ParentID,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		}
	}

	return &v1.ListCategoriesResponse{
		Categories: v1Category,
		Pagination: &v1.PaginationResponse{
			PageSize:  req.Pagination.GetPageSize(),
			PageToken: nextPageToken,
			AfterPage: afterPage,
		},
	}, nil
}

func (api *serverAPI) DeleteCategory(context.Context, *v1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	panic("implement me")
}
