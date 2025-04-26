package product

import (
	"context"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	v1 "github.com/DENFNC/Zappy/catalog_service/proto/gen/v1"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc"
)

type Product interface {
	Create(
		ctx context.Context,
		name string,
		desc, currency string,
		price pgtype.Numeric,
		categoryIDs []string,
	) (string, error)
	Get(
		ctx context.Context,
		productID string,
	) (*models.Product, error)
	List(
		ctx context.Context,
		page, pageSize int32,
		query string,
		categoryIDs []string,
	) ([]*models.Product, int32, error)
	Update(
		ctx context.Context,
		productID, name, description string,
		price pgtype.Numeric,
		currency string,
		categoryIDs []string,
		isActive *bool,
	) (*models.Product, error)
	Delete(
		ctx context.Context,
		productID string,
	) error
}

type serverAPI struct {
	v1.UnimplementedProductServiceServer
	svc Product
}

func New(svc Product) *serverAPI {
	return &serverAPI{
		svc: svc,
	}
}

func (api *serverAPI) Register(grpc *grpc.Server) {
	v1.RegisterProductServiceServer(grpc, api)
}

func (api *serverAPI) CreateProduct(
	ctx context.Context,
	req *v1.CreateProductRequest,
) (*v1.CreateProductResponse, error) {
	panic("implement me")
}

func (api *serverAPI) GetProduct(
	ctx context.Context,
	req *v1.GetProductRequest,
) (*v1.GetProductResponse, error) {
	panic("implement me")
}

func (api *serverAPI) ListProducts(
	ctx context.Context,
	req *v1.ListProductsRequest,
) (*v1.ListProductsResponse, error) {
	panic("implement me")
}

func (api *serverAPI) UpdateProduct(
	ctx context.Context,
	req *v1.UpdateProductRequest,
) (*v1.UpdateProductResponse, error) {
	panic("implement me")
}
