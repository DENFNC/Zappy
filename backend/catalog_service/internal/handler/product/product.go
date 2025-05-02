package product

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

type Product interface {
	Create(
		ctx context.Context,
		name, desc string,
		categoryIDs []string,
		price int64,
	) (string, error)
	Get(
		ctx context.Context,
		productID string,
	) (*models.Product, error)
	List(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
	) ([]models.Product, string, error)
	Update(
		ctx context.Context,
		productID, name, description string,
		price int64,
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
	productID, err := api.svc.Create(
		ctx,
		req.GetName(),
		req.GetDescription(),
		req.GetCategoryIds(),
		req.GetPriceCents(),
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.CreateProductResponse{
		ProductId: &v1.ResourceID{
			Id: productID,
		},
	}, nil
}

func (api *serverAPI) GetProduct(
	ctx context.Context,
	req *v1.GetProductRequest,
) (*v1.GetProductResponse, error) {
	product, err := api.svc.Get(ctx, req.GetProductId().GetId())
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.GetProductResponse{
		Product: &v1.Product{
			Id:          product.ProductID,
			Name:        product.ProductName,
			Description: product.Description,
			PriceCents:  product.Price,
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		},
	}, nil
}

func (api *serverAPI) ListProducts(
	ctx context.Context,
	req *v1.ListProductsRequest,
) (*v1.ListProductsResponse, error) {
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

	v1Products := make([]*v1.Product, len(items))
	for i, item := range items {
		v1Products[i] = &v1.Product{
			Id:          item.ProductID,
			Name:        item.ProductName,
			Description: item.Description,
			PriceCents:  item.Price,
			CreatedAt:   timestamppb.New(item.CreatedAt),
			UpdatedAt:   timestamppb.New(item.UpdatedAt),
		}
	}

	return &v1.ListProductsResponse{
		Products: v1Products,
		Pagination: &v1.PaginationResponse{
			PageSize:  req.Pagination.GetPageSize(),
			PageToken: nextPageToken,
			AfterPage: afterPage,
		},
	}, nil
}

func (api *serverAPI) UpdateProduct(
	ctx context.Context,
	req *v1.UpdateProductRequest,
) (*v1.UpdateProductResponse, error) {
	panic("implement me")
}

func (api *serverAPI) DeleteProduct(
	ctx context.Context,
	req *v1.DeleteProductRequest,
) (*emptypb.Empty, error) {
	err := api.svc.Delete(ctx, req.GetProductId())
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &emptypb.Empty{}, nil
}
