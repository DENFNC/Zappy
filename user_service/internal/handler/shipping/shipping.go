package shipping

import (
	"context"

	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Shipping interface{}

type serverAPI struct {
	v1.UnimplementedShippingAddressServiceServer
	service Shipping
}

func Register(grpc *grpc.Server, shipping Shipping) {
	v1.RegisterShippingAddressServiceServer(grpc, &serverAPI{service: shipping})
}

func (sa *serverAPI) CreateShippingAddress(context.Context, *v1.CreateShippingAddressRequest) (*v1.ShippingAddress, error) {
	panic("implement me!")
}

func (sa *serverAPI) DeleteShippingAddress(context.Context, *v1.DeleteShippingAddressRequest) (*emptypb.Empty, error) {
	panic("implement me!")
}

func (sa *serverAPI) GetShippingAddress(context.Context, *v1.GetShippingAddressRequest) (*v1.ShippingAddress, error) {
	panic("implement me!")
}

func (sa *serverAPI) ListShippingAddresses(context.Context, *v1.ListShippingAddressesRequest) (*v1.ListShippingAddressesResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) UpdateShippingAddress(context.Context, *v1.UpdateShippingAddressRequest) (*v1.ShippingAddress, error) {
	panic("implement me!")
}
