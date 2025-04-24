package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Profile interface {
	Create(ctx context.Context, authUserID string, firstName, lastName string) (string, error)
	Delete(ctx context.Context, profileID string) (string, error)
	GetByID(ctx context.Context, profileID string) (*models.Profile, error)
	List(ctx context.Context, params []any) ([]any, string, error)
	Update(ctx context.Context, profileID string, firstName, lastName string) (string, error)
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

func (sa *serverAPI) CreateProfile(ctx context.Context, req *v1.CreateProfileRequest) (*v1.ProfileIDResponse, error) {
	profileID, err := sa.service.Create(
		ctx,
		req.GetProfile().GetAuthUserId(),
		req.GetProfile().GetName().GetFirstName(),
		req.GetProfile().GetName().GetLastName(),
	)

	fmt.Println(err)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &v1.ProfileIDResponse{
		ProfileId: profileID,
	}, nil
}

func (sa *serverAPI) DeleteProfile(ctx context.Context, req *v1.DeleteProfileRequest) (*v1.ProfileIDResponse, error) {
	profileID, err := sa.service.Delete(
		ctx,
		req.GetProfileId(),
	)

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &v1.ProfileIDResponse{
		ProfileId: profileID,
	}, nil
}

func (sa *serverAPI) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.Profile, error) {
	profile, err := sa.service.GetByID(ctx, req.GetProfileId())

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &v1.Profile{
		ProfileId:  profile.ProfileID,
		AuthUserId: profile.AuthUserID,
		Name: &v1.FullName{
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
		},
		CreatedAt: timestamppb.New(profile.CreatedAt),
		UpdatedAt: timestamppb.New(profile.UpdatedAt),
	}, nil
}

func (sa *serverAPI) ListProfiles(ctx context.Context, req *v1.ListProfilesRequest) (*v1.ListProfilesResponse, error) {
	panic("implement me")
}

func (sa *serverAPI) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (*v1.ProfileIDResponse, error) {
	profileID, err := sa.service.Update(
		ctx,
		req.GetProfileId(),
		req.GetProfile().GetFirstName(),
		req.GetProfile().GetLastName(),
	)

	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &v1.ProfileIDResponse{
		ProfileId: profileID,
	}, nil
}
