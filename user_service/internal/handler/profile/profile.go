package profile

import (
	"context"

	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Profile interface{}

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
func (sa *serverAPI) CreateProfile(context.Context, *v1.CreateProfileRequest) (*v1.Profile, error) {
	panic("implement me!")
}

func (sa *serverAPI) DeleteProfile(context.Context, *v1.DeleteProfileRequest) (*emptypb.Empty, error) {
	panic("implement me!")
}

func (sa *serverAPI) GetProfile(context.Context, *v1.GetProfileRequest) (*v1.Profile, error) {
	panic("implement me!")
}

func (sa *serverAPI) ListProfiles(context.Context, *v1.ListProfilesRequest) (*v1.ListProfilesResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) UpdateProfile(context.Context, *v1.UpdateProfileRequest) (*v1.Profile, error) {
	panic("implement me!")
}
