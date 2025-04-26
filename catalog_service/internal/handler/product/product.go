package product

import (
	"context"

	v1 "github.com/DENFNC/Zappy/catalog_service/proto/gen/v1"
	"google.golang.org/grpc"
)

type Product interface{}

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
