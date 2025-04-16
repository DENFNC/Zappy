package profile

import (
	"context"
	"fmt"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Profile interface {
	Create(ctx context.Context, authUserID uint64, firstName, lastName string) (uint64, error)
	Delete(ctx context.Context, profileID int) (uint64, error)
	GetByID(ctx context.Context, profileID int) (*models.Profile, error)
	List(ctx context.Context) ([]*models.Profile, error)
	Update(ctx context.Context, profileID int, firstName, lastName, phone string) (uint64, error)
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
		req.GetProfile().AuthUserId,
		req.GetProfile().Name.FirstName,
		req.GetProfile().Name.LastName,
	)

	fmt.Println(err)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &v1.ProfileIDResponse{
		ProfileId: profileID,
		Message:   "Пользователь создан",
	}, nil
}

func (sa *serverAPI) DeleteProfile(ctx context.Context, req *v1.DeleteProfileRequest) (*v1.ProfileIDResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.Profile, error) {
	panic("implement me!")
}

func (sa *serverAPI) ListProfiles(ctx context.Context, req *v1.ListProfilesRequest) (*v1.ListProfilesResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (*v1.ProfileIDResponse, error) {
	panic("implement me!")
}
