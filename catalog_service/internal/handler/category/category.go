package category

import (
	"context"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	v1 "github.com/DENFNC/Zappy/catalog_service/proto/gen/v1"
	"google.golang.org/grpc"
)

type Category interface {
	CreateCategory(
		ctx context.Context,
		name, parentID string,
	) (*models.Category, error)
	ListCategories(
		ctx context.Context,
		page, pageSize int32,
	) ([]models.Category, int32, error)
	DeleteCategory(
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
	panic("implement me")
}

func (api *serverAPI) ListCategories(
	ctx context.Context,
	req *v1.ListCategoriesRequest,
) (*v1.ListCategoriesResponse, error) {
	panic("implement me")
}
