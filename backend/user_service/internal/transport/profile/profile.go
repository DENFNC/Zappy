package profile

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/DENFNC/Zappy/user_service/proto/gen/go/common/v1"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/go/profile/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Profile interface {
	Create(
		ctx context.Context,
		profile *models.Profile,
	) (string, error)
	Delete(
		ctx context.Context,
		profileID string,
	) (string, error)
	GetByID(
		ctx context.Context,
		profileID string,
	) (*models.Profile, error)
	List(
		ctx context.Context,
		params []any,
	) ([]any, string, error)
	Update(
		ctx context.Context,
		profileID string,
		firstName, lastName string,
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

func (sa *serverAPI) Register(grpc *grpc.Server) {
	v1.RegisterUserProfileServiceServer(grpc, sa)
}

func (sa *serverAPI) CreateProfile(ctx context.Context, req *v1.CreateProfileRequest) (*v1.CreateProfileResponse, error) {
	profileID, err := sa.service.Create(
		ctx,
		&models.Profile{
			FirstName: req.Profile.Name.GetFirstName(),
			LastName:  req.Profile.Name.LastName,
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &v1.CreateProfileResponse{
		ProfileId: &common.ResourceID{
			Id: profileID,
		},
	}, nil
}

func (sa *serverAPI) DeleteProfile(ctx context.Context, req *v1.DeleteProfileRequest) (*v1.DeleteProfileResponse, error) {
	profileID, err := sa.service.Delete(
		ctx,
		req.ProfileId.GetId(),
	)

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &v1.DeleteProfileResponse{
		ProfileId: &common.ResourceID{
			Id: profileID,
		},
	}, nil
}

func (sa *serverAPI) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	profile, err := sa.service.GetByID(ctx, req.ProfileId.GetId())

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
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

func (sa *serverAPI) ListProfiles(ctx context.Context, req *v1.ListProfilesRequest) (*v1.ListProfilesResponse, error) {
	profiles, _, err := sa.service.List(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	v1Profiles := make([]*v1.Profile, len(profiles))
	for i, p := range profiles {
		profile, ok := p.(*models.Profile)
		if !ok {
			return nil, status.Error(codes.Internal, "Failed to cast profile")
		}
		v1Profiles[i] = &v1.Profile{
			ProfileId:  profile.ProfileID,
			AuthUserId: profile.AuthUserID,
			Name: &v1.FullName{
				FirstName: profile.FirstName,
				LastName:  profile.LastName,
			},
			CreatedAt: timestamppb.New(profile.CreatedAt),
			UpdatedAt: timestamppb.New(profile.UpdatedAt),
		}
	}

	return &v1.ListProfilesResponse{
		Profiles: v1Profiles,
	}, nil
}

func (sa *serverAPI) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (*v1.UpdateProfileResponse, error) {
	profileID, err := sa.service.Update(
		ctx,
		req.ProfileId.GetId(),
		req.GetProfile().GetFirstName(),
		req.GetProfile().GetLastName(),
	)

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &v1.UpdateProfileResponse{
		ProfileId: &common.ResourceID{
			Id: profileID,
		},
	}, nil
}
