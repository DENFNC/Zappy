package productimage

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/dbutils"
	errpkg "github.com/DENFNC/Zappy/catalog_service/internal/utils/errors"
	"github.com/DENFNC/Zappy/catalog_service/proto/gen/go/common/v1"
	v1 "github.com/DENFNC/Zappy/catalog_service/proto/gen/go/product_image/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var statusMap = map[string]v1.Status{
	"pending": v1.Status_STATUS_PENDING,
	"success": v1.Status_STATUS_SUCCESS,
	"failed":  v1.Status_STATUS_FAILED,
}

type ProductImage interface {
	GetByID(
		ctx context.Context,
		imageID string,
	) (*models.ProductImage, error)
	ListProductImages(
		ctx context.Context,
		pageSize uint32,
		pageToken, productID string,
	) ([]models.ProductImage, string, error)
	DeleteProductImage(
		ctx context.Context,
		imageID string,
	) error
	GetUploadURL(
		ctx context.Context,
		bucket, key, contentType string,
		productID, alt string,
	) (string, error)
	GetUploadStatus(
		ctx context.Context,
		key string,
	) (string, error)
}

type serverAPI struct {
	v1.UnimplementedProductImageServiceServer
	svc        ProductImage
	bucketName string
}

func New(svc ProductImage, bucketName string) *serverAPI {
	return &serverAPI{
		svc:        svc,
		bucketName: bucketName,
	}
}

func (api *serverAPI) GRPCRegister(grpc *grpc.Server) {
	v1.RegisterProductImageServiceServer(grpc, api)
}

func (api *serverAPI) HTTPRegister(
	ctx context.Context,
	mux *runtime.ServeMux,
) {
	v1.RegisterProductImageServiceHandlerServer(ctx, mux, api)
}

func (api *serverAPI) GetProductImage(
	ctx context.Context,
	req *v1.GetProductImageRequest,
) (*v1.GetProductImageResponse, error) {
	item, err := api.svc.GetByID(
		ctx,
		req.GetImageId(),
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.GetProductImageResponse{
		ProductImage: &v1.ProductImage{
			ImageId:   item.ImageID,
			ProductId: item.ProductID,
			Url:       item.URL,
			Alt:       item.ALT,
			ObjectKey: item.ObjectKey,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		},
	}, nil
}

func (api *serverAPI) ListProductImages(
	ctx context.Context,
	req *v1.ListProductImagesRequest,
) (*v1.ListProductImagesResponse, error) {
	afterPage := true
	items, nextPageToken, err := api.svc.ListProductImages(
		ctx,
		req.Pagination.GetPageSize(),
		req.Pagination.GetPageToken(),
		req.GetProductId(),
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}
	if nextPageToken == "" {
		afterPage = false
	}

	v1Items := make([]*v1.ProductImage, len(items))
	for i, item := range items {
		v1Items[i] = &v1.ProductImage{
			ImageId:   item.ImageID,
			ProductId: item.ProductID,
			Url:       item.URL,
			Alt:       item.ALT,
			ObjectKey: item.ObjectKey,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		}
	}

	return &v1.ListProductImagesResponse{
		ProductImage: v1Items,
		Pagination: &common.PaginationResponse{
			PageSize:  uint32(len(v1Items)),
			PageToken: nextPageToken,
			AfterPage: afterPage,
		},
	}, nil
}

func (api *serverAPI) DeleteProductImage(
	ctx context.Context,
	req *v1.DeleteProductImageRequest,
) (*v1.DeleteProductImageResponse, error) {
	if err := api.svc.DeleteProductImage(
		ctx,
		req.GetImageId(),
	); err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.DeleteProductImageResponse{}, nil
}

func (api *serverAPI) UpdateProductImage(
	ctx context.Context,
	req *v1.UpdateProductImageRequest,
) (*v1.UpdateProductImageResponse, error) {
	panic("implement me!")
}

func (api *serverAPI) UploadFileURL(
	ctx context.Context,
	req *v1.UploadFileURLRequest,
) (*v1.UploadFileURLResponse, error) {
	uid := dbutils.NewUUIDV7().String()

	var builder strings.Builder
	builder.WriteString(uid)
	builder.WriteString(filepath.Ext(req.GetFilename()))
	key := builder.String()

	url, err := api.svc.GetUploadURL(ctx, api.bucketName, key,
		req.GetContentType(),
		req.GetProductId(),
		req.GetAlt(),
	)

	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.UploadFileURLResponse{
		UploadUrl: url,
		ObjectKey: key,
	}, nil
}

func (api *serverAPI) UploadStatus(
	ctx context.Context,
	req *v1.UploadStatusRequest,
) (*v1.UploadStatusResponse, error) {
	status, _ := api.svc.GetUploadStatus(
		ctx,
		req.GetKey(),
	)

	return &v1.UploadStatusResponse{
		Key:    req.GetKey(),
		Status: v1.Status(statusMap[status]),
	}, nil
}
