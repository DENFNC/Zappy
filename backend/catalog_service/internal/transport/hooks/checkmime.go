package hooks

import (
	"context"
	"fmt"

	hooks "github.com/DENFNC/Zappy/catalog_service/proto/gen/go/object/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type WebHook interface {
	CheckMIME(data []byte) error
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
	fmt.Println(req.String())

	return &hooks.WebHookServiceCheckMimeObjectStorageResponse{}, nil
}
