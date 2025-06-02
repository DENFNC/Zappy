package auth

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/auth_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/auth_service/internal/utils/errors"
	v1 "github.com/DENFNC/Zappy/auth_service/proto/gen/go/auth/v1"
	"github.com/DENFNC/Zappy/auth_service/proto/gen/go/common/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	UserLogin(
		ctx context.Context,
		identifier string,
		password string,
	) (string, error)
	UserRegister(
		ctx context.Context,
		user *models.User,
	) (string, uint64, error)
	UpdateRefreshToken(
		ctx context.Context,
		token string,
	) (string, error)
	ListUsers(
		ctx context.Context,
		pageSize uint,
		pageToken string,
	) ([]*models.User, string, error)
}

type serverAPI struct {
	v1.UnimplementedAuthServiceServer
	svc Auth
}

func New(svc Auth) *serverAPI {
	return &serverAPI{
		svc: svc,
	}
}

func (api *serverAPI) GRPCRegister(gRPC *grpc.Server) {
	v1.RegisterAuthServiceServer(gRPC, api)
}

func (api *serverAPI) HTTPRegister(
	ctx context.Context,
	mux *runtime.ServeMux,
) {
	v1.RegisterAuthServiceHandlerServer(ctx, mux, api)
}

func (api *serverAPI) Login(
	ctx context.Context,
	req *v1.LoginRequest,
) (*v1.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			errpkg.ErrInvalidArgument.Message,
		)
	}

	var authIdentifier string
	switch identifier := req.AuthType.(type) {
	case *v1.LoginRequest_Email:
		authIdentifier = identifier.Email.GetEmail()
	case *v1.LoginRequest_Username:
		authIdentifier = identifier.Username.GetUsername()
	default:
		return nil, status.Error(
			codes.InvalidArgument,
			"unsupported auth type",
		)
	}

	token, err := api.svc.UserLogin(ctx, authIdentifier, req.GetPassword().Password)
	if err != nil {
		var appErr *errpkg.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case errpkg.ErrInvalidCredentials.Code:
				return nil, status.Error(
					codes.Unauthenticated,
					appErr.Message,
				)
			case errpkg.ErrUserNotFound.Code:
				return nil, status.Error(
					codes.NotFound,
					appErr.Message,
				)
			default:
				return nil, status.Error(
					codes.Internal,
					errpkg.ErrInternal.Message,
				)
			}
		}
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.LoginResponse{
		Token: token,
	}, nil
}

func (api *serverAPI) Register(
	ctx context.Context,
	req *v1.RegisterRequest,
) (*v1.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			errpkg.ErrInvalidArgument.Message,
		)
	}

	token, userID, err := api.svc.UserRegister(
		ctx,
		&models.User{
			Username: req.GetUser().Username.Username,
			Email:    req.GetUser().Email.Email,
			Password: []byte(req.GetPassword().Password),
		},
	)
	if err != nil {
		var appErr *errpkg.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case errpkg.ErrUserAlreadyExists.Code:
				return nil, status.Error(
					codes.AlreadyExists,
					appErr.Message,
				)
			default:
				return nil, status.Error(
					codes.Internal,
					errpkg.ErrInternal.Message,
				)
			}
		}
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.RegisterResponse{
		Token:  token,
		UserId: userID,
	}, nil
}

func (api *serverAPI) ListUsers(
	ctx context.Context,
	req *v1.ListUsersRequest,
) (*v1.ListUsersResponse, error) {
	afterPage := true
	users, nextPageToken, err := api.svc.ListUsers(
		ctx,
		uint(req.Pagination.GetPageSize()),
		req.Pagination.GetPageToken(),
	)
	if err != nil {
		var appErr *errpkg.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case errpkg.ErrInvalidArgument.Code:
				return nil, status.Error(
					codes.InvalidArgument,
					appErr.Message,
				)
			default:
				return nil, status.Error(
					codes.Internal,
					errpkg.ErrInternal.Message,
				)
			}
		}
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}
	if nextPageToken == "" {
		afterPage = false
	}

	protoUsers := make([]*v1.UserDefault, len(users))
	for i, user := range users {
		protoUsers[i] = &v1.UserDefault{
			Username: &v1.Username{Username: user.Username},
			Email:    &v1.Email{Email: user.Email},
		}
	}

	return &v1.ListUsersResponse{
		Users: protoUsers,
		Pagination: &common.PaginationResponse{
			PageSize:  uint32(len(protoUsers)),
			PageToken: nextPageToken,
			AfterPage: afterPage,
		},
	}, nil
}

func (api *serverAPI) Refresh(
	ctx context.Context,
	req *v1.RefreshRequest,
) (*v1.RefreshResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			errpkg.ErrInvalidArgument.Message,
		)
	}

	newToken, err := api.svc.UpdateRefreshToken(ctx, req.GetToken())
	if err != nil {
		var appErr *errpkg.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case errpkg.ErrInvalidToken.Code:
				return nil, status.Error(
					codes.Unauthenticated,
					appErr.Message,
				)
			default:
				return nil, status.Error(
					codes.Internal,
					errpkg.ErrInternal.Message,
				)
			}
		}
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.RefreshResponse{
		Token: newToken,
	}, nil
}
