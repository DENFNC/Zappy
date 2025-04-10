package authgrpc

import (
	"context"
	"fmt"

	v1 "github.com/DENFNC/Zappy/proto/gen/v1"
	"google.golang.org/grpc"
)

type Auth interface {
	Login(
		ctx context.Context,
		authType string,
		password string,
	) (string, error)
	Register(
		ctx context.Context,
		username string,
		email string,
		password string,
	) (string, uint64, error)
	Refresh(
		ctx context.Context,
		token string,
	) (string, error)
}

type serverAPI struct {
	v1.UnimplementedAuthServer
	// log  *slog.Logger
	auth Auth
}

func ServRegister(gRPC *grpc.Server, auth Auth) {
	v1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (sa *serverAPI) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	panic("implement me")
}

func (sa *serverAPI) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	token, userID, err := sa.auth.Register(
		ctx,
		req.GetUser().Username.Username,
		req.GetUser().Email.Email,
		req.GetPassword().Password,
	)

	if err != nil {
		fmt.Println(err)
	}

	return &v1.RegisterResponse{
		Token:  token,
		UserId: userID,
	}, nil
}

func (sa *serverAPI) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	panic("implement me")
}
