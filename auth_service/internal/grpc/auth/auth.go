package authgrpc

import (
	"context"

	v1 "github.com/DENFNC/Zappy/proto/gen/v1"
	"google.golang.org/grpc"
)

type Auth interface {
	Login(
		ctx context.Context,
		authType string,
		password string,
	)
	Register(
		ctx context.Context,
		username string,
		email string,
		password string,
	)
	Refresh(
		ctx context.Context,
		token string,
	)
}

type serverAPI struct {
	v1.UnimplementedAuthServer
	auth Auth
}

func ServRegister(gRPC *grpc.Server, auth Auth) {
	v1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (sa *serverAPI) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	panic("implement me")
}

func (sa *serverAPI) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	panic("implement me")
}

func (sa *serverAPI) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	panic("implement me")
}
