package authgrpc

import (
	"context"
	"errors"

	errpkg "github.com/DENFNC/Zappy/internal/errors"
	v1 "github.com/DENFNC/Zappy/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Auth – интерфейс, который предоставляет методы аутентификации.
type Auth interface {
	Login(ctx context.Context, authType string, password string) (string, error)
	Register(ctx context.Context, username string, email string, password string) (string, uint64, error)
	Refresh(ctx context.Context, token string) (string, error)
}

type serverAPI struct {
	v1.UnimplementedAuthServer
	auth Auth
}

// ServRegister регистрирует сервер gRPC.
func ServRegister(gRPC *grpc.Server, auth Auth) {
	v1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

// Login реализует метод входа через gRPC.
func (sa *serverAPI) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrInvalidArgument.Message)
	}

	var authIdentifier string
	switch identifier := req.AuthType.(type) {
	case *v1.LoginRequest_Email:
		authIdentifier = identifier.Email.GetEmail()
	case *v1.LoginRequest_Username:
		authIdentifier = identifier.Username.GetUsername()
	default:
		return nil, status.Error(codes.InvalidArgument, "unsupported auth type")
	}

	token, err := sa.auth.Login(ctx, authIdentifier, req.GetPassword().Password)
	if err != nil {
		var appErr *errpkg.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case "INVALID_CREDENTIALS":
				return nil, status.Error(codes.Unauthenticated, appErr.Message)
			default:
				return nil, status.Error(codes.Internal, appErr.Message)
			}
		}
		return nil, status.Error(codes.Internal, errpkg.ErrInternalServer.Message)
	}

	return &v1.LoginResponse{
		Token: token,
	}, nil
}

func (sa *serverAPI) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrInvalidArgument.Message)
	}

	token, userID, err := sa.auth.Register(
		ctx,
		req.GetUser().Username.Username,
		req.GetUser().Email.Email,
		req.GetPassword().Password,
	)
	if err != nil {
		var appErr *errpkg.AppError
		if errors.As(err, &appErr) && appErr.Code == "INVALID_CREDENTIALS" {
			return nil, status.Error(codes.Unauthenticated, appErr.Message)
		}
		return nil, status.Error(codes.Internal, errpkg.ErrInternalServer.Message)
	}

	return &v1.RegisterResponse{
		Token:  token,
		UserId: userID,
	}, nil
}

func (sa *serverAPI) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrInvalidArgument.Message)
	}

	newToken, err := sa.auth.Refresh(ctx, req.GetToken())
	if err != nil {
		var appErr *errpkg.AppError
		if errors.As(err, &appErr) && appErr.Code == "INVALID_TOKEN" {
			return nil, status.Error(codes.Unauthenticated, appErr.Message)
		}
		return nil, status.Error(codes.Internal, errpkg.ErrInternalServer.Message)
	}

	return &v1.RefreshResponse{
		Token: newToken,
	}, nil
}
