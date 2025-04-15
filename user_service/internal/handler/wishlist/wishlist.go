package wishlist

import (
	"context"

	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Wishlist interface{}

type serverAPI struct {
	v1.UnimplementedWishlistItemServiceServer
	service Wishlist
}

func Register(grpc *grpc.Server, wishlist Wishlist) {
	v1.RegisterWishlistItemServiceServer(grpc, &serverAPI{service: wishlist})
}

func (sa *serverAPI) CreateWishlistItem(context.Context, *v1.CreateWishlistItemRequest) (*v1.WishlistItem, error) {
	panic("implement me!")
}

func (sa *serverAPI) DeleteWishlistItem(context.Context, *v1.DeleteWishlistItemRequest) (*emptypb.Empty, error) {
	panic("implement me!")
}

func (sa *serverAPI) GetWishlistItem(context.Context, *v1.GetWishlistItemRequest) (*v1.WishlistItem, error) {
	panic("implement me!")
}

func (sa *serverAPI) ListWishlistItems(*v1.ListWishlistItemsRequest, grpc.ServerStreamingServer[v1.WishlistItem]) error {
	panic("implement me!")
}

func (sa *serverAPI) UpdateWishlistItem(context.Context, *v1.UpdateWishlistItemRequest) (*v1.WishlistItem, error) {
	panic("implement me!")
}
