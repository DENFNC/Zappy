package profile

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/domain/models"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
)

type Profile interface {
	GetProfile(
		id uint64,
	) (models.Profile, error)
	UpdateProfile(
		firstName, lastName string,
		avatarURL []string,
	) (string, error)
	DeleteProfile(
		id uint64,
	) (string, error)
}

type serverAPI struct {
	v1.UnimplementedProfileServer
	prf Profile
}

func ProfileRegister(grpc *grpc.Server, prf Profile) {
	v1.RegisterProfileServer(grpc, &serverAPI{prf: prf})
}

func (sa *serverAPI) DeleteProfile(context.Context, *v1.DeleteProfileRequest) (*v1.DeleteProfileResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) GetProfile(context.Context, *v1.ProfileRequest) (*v1.ProfileResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) UpdateProfile(context.Context, *v1.UpdateProfileRequest) (*v1.UpdateProfileResponse, error) {
	panic("implement me!")
}
