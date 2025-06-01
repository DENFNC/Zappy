package profile

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
	"github.com/DENFNC/Zappy/user_service/proto/gen/go/common/v1"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/go/profile/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Profile interface {
	CreateProfile(
		ctx context.Context,
		profile *models.Profile,
	) (string, error)
	DeleteProfile(
		ctx context.Context,
		profileID string,
	) (string, error)
	ProfileGetByID(
		ctx context.Context,
		profileID string,
	) (*models.Profile, error)
	ListProfiles(
		ctx context.Context,
		pageSize uint,
		pageToken string,
	) ([]*models.Profile, string, error)
	UpdateProfile(
		ctx context.Context,
		profileID string,
		profile *models.Profile,
	) (string, error)
}

type serverAPI struct {
	v1.UnimplementedUserProfileServiceServer
	service Profile
}

func New(service Profile) *serverAPI {
	return &serverAPI{
		service: service,
	}
}

func (api *serverAPI) GRPCRegister(grpc *grpc.Server) {
	v1.RegisterUserProfileServiceServer(grpc, api)
}

func (api *serverAPI) HTTPRegister(
	ctx context.Context,
	mux *runtime.ServeMux,
) {
	v1.RegisterUserProfileServiceHandlerServer(ctx, mux, api)
}

func (api *serverAPI) CreateProfile(ctx context.Context, req *v1.CreateProfileRequest) (*v1.CreateProfileResponse, error) {
	profileID, err := api.service.CreateProfile(
		ctx,
		&models.Profile{
			FirstName: req.Profile.Name.GetFirstName(),
			LastName:  req.Profile.Name.LastName,
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, errpkg.ErrInternal.Message)
	}

	return &v1.CreateProfileResponse{
		ProfileId: &common.ResourceID{
			Id: profileID,
		},
	}, nil
}

func (api *serverAPI) DeleteProfile(ctx context.Context, req *v1.DeleteProfileRequest) (*v1.DeleteProfileResponse, error) {
	profileID, err := api.service.DeleteProfile(
		ctx,
		req.ProfileId.GetId(),
	)

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, errpkg.ErrNotFound.Message)
		}
		return nil, status.Error(codes.Internal, errpkg.ErrInternal.Message)
	}

	return &v1.DeleteProfileResponse{
		ProfileId: &common.ResourceID{
			Id: profileID,
		},
	}, nil
}

func (api *serverAPI) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	profile, err := api.service.ProfileGetByID(ctx, req.ProfileId.GetId())

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, errpkg.ErrNotFound.Message)
		}
		return nil, status.Error(codes.Internal, errpkg.ErrInternal.Message)
	}

	return &v1.GetProfileResponse{
		Profile: &v1.Profile{
			ProfileId:  profile.ProfileID,
			AuthUserId: profile.AuthUserID,
			Name: &v1.FullName{
				FirstName: profile.FirstName,
				LastName:  profile.LastName,
			},
			CreatedAt: timestamppb.New(profile.CreatedAt),
			UpdatedAt: timestamppb.New(profile.UpdatedAt),
		},
	}, nil
}

func (api *serverAPI) ListProfiles(ctx context.Context, req *v1.ListProfilesRequest) (*v1.ListProfilesResponse, error) {
	afterPage := true
	items, pageToken, err := api.service.ListProfiles(
		ctx,
		uint(req.Pagination.GetPageSize()),
		req.Pagination.GetPageToken(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, errpkg.ErrInternal.Message)
	}
	if pageToken == "" {
		afterPage = false
	}

	v1Profiles := make([]*v1.Profile, len(items))
	for i, item := range items {
		v1Profiles[i] = &v1.Profile{
			ProfileId:  item.ProfileID,
			AuthUserId: item.AuthUserID,
			Name: &v1.FullName{
				FirstName: item.FirstName,
				LastName:  item.LastName,
			},
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		}
	}

	return &v1.ListProfilesResponse{
		Profiles: v1Profiles,
		Pagination: &common.PaginationResponse{
			PageSize:  uint32(len(items)),
			PageToken: pageToken,
			AfterPage: afterPage,
		},
	}, nil
}

func (api *serverAPI) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (*v1.UpdateProfileResponse, error) {
	profileID, err := api.service.UpdateProfile(
		ctx,
		req.ProfileId.GetId(),
		&models.Profile{
			FirstName: req.Profile.GetFirstName(),
			LastName:  req.Profile.GetLastName(),
		},
	)

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, errpkg.ErrNotFound.Message)
		}
		return nil, status.Error(codes.Internal, errpkg.ErrInternal.Message)
	}

	return &v1.UpdateProfileResponse{
		ProfileId: &common.ResourceID{
			Id: profileID,
		},
	}, nil
}
