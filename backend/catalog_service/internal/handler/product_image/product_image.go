package productimage

import (
	"context"

	errpkg "github.com/DENFNC/Zappy/catalog_service/internal/utils/errors"
	v1 "github.com/DENFNC/Zappy/catalog_service/proto/gen/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductImage interface {
	GetUploadURL(
		ctx context.Context,
		bucket, key, contentType string,
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

func (api *serverAPI) UploadFileURL(
	ctx context.Context,
	req *v1.UploadURLRequest,
) (*v1.UploadURLResponse, error) {
	url, err := api.svc.GetUploadURL(
		ctx,
		api.bucketName,
		req.GetFilename(),
		req.GetContentType(),
	)

	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.UploadURLResponse{
		UploadUrl: url,
		ObjectKey: req.GetFilename(),
	}, nil
}
