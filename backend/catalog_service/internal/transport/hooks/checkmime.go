package hooks

import (
	"context"

	hooks "github.com/DENFNC/Zappy/catalog_service/proto/gen/go/object/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type WebHook interface {
	CheckMime(
		ctx context.Context,
		bucket, key string,
		byteRange string,
	) error
}

type serverAPI struct {
	hooks.UnimplementedWebHookServiceServer
	svc WebHook
}

func New(svc WebHook) *serverAPI {
	return &serverAPI{
		svc: svc,
	}
}

func (api *serverAPI) GRPCRegister(grpc *grpc.Server) {
	hooks.RegisterWebHookServiceServer(grpc, api)
}

func (api *serverAPI) HTTPRegister(
	ctx context.Context,
	mux *runtime.ServeMux,
) {
	hooks.RegisterWebHookServiceHandlerServer(ctx, mux, api)
}

func (api *serverAPI) CheckMimeObjectStorage(
	ctx context.Context,
	req *hooks.WebHookServiceCheckMimeObjectStorageRequest,
) (*hooks.WebHookServiceCheckMimeObjectStorageResponse, error) {
	group, ctx := errgroup.WithContext(ctx)
	hookRecords := req.GetRecords()

	for _, rec := range hookRecords {
		group.Go(func() error {
			return api.svc.CheckMime(
				ctx,
				rec.GetS3().Bucket.Name,
				rec.GetS3().Object.Key,
				"bytes=0-511",
			)
		})
	}
	if err := group.Wait(); err != nil {
		return nil, err
	}

	return &hooks.WebHookServiceCheckMimeObjectStorageResponse{}, nil
}
