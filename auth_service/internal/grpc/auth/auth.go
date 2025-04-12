package authgrpc

import (
	"context"
	"errors"

	authservice "github.com/DENFNC/Zappy/internal/service/auth"
	v1 "github.com/DENFNC/Zappy/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	auth Auth
}

func ServRegister(gRPC *grpc.Server, auth Auth) {
	v1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (sa *serverAPI) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	var authIdentifier string

	switch identifier := req.AuthType.(type) {
	case *v1.LoginRequest_Email:
		authIdentifier = identifier.Email.GetEmail()
	case *v1.LoginRequest_Username:
		authIdentifier = identifier.Username.GetUsername()
	default:
		return nil, status.Error(codes.InvalidArgument, "Unsupported auth type")
	}

	token, err := sa.auth.Login(
		ctx,
		authIdentifier,
		req.GetPassword().Password,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	return &v1.LoginResponse{
		Token: token,
	}, nil
}

func (sa *serverAPI) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	token, userID, err := sa.auth.Register(
		ctx,
		req.GetUser().Username.Username,
		req.GetUser().Email.Email,
		req.GetPassword().Password,
	)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		// На остальные ошибки возвращаем внутреннюю ошибку сервера.
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &v1.RegisterResponse{
		Token:  token,
		UserId: userID,
	}, nil
}

func (sa *serverAPI) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	panic("implement me")
}
