package wishlist

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/DENFNC/Zappy/user_service/proto/gen/go/common/v1"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/go/wishlist/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Wishlist interface {
	CreateItem(ctx context.Context, profileID, productID string) (string, error)
	GetItem(ctx context.Context, itemID string) (*models.WishlistItem, error)
	UpdateItem(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error)
	DeleteItem(ctx context.Context, itemID string) error
	ListItems(ctx context.Context, profileID string) ([]*models.WishlistItem, error)
}

type serverAPI struct {
	v1.UnimplementedWishlistItemServiceServer
	service Wishlist
}

func New(service Wishlist) *serverAPI {
	return &serverAPI{
		service: service,
	}
}

func (sa *serverAPI) Register(grpc *grpc.Server) {
	v1.RegisterWishlistItemServiceServer(grpc, sa)
}

func (sa *serverAPI) CreateWishlistItem(ctx context.Context, req *v1.CreateWishlistItemRequest) (*v1.CreateWishlistItemResponse, error) {
	itemID, err := sa.service.CreateItem(
		ctx,
		req.GetWishlistItem().GetProfileId(),
		req.GetWishlistItem().GetProductId(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create wishlist item")
	}

	return &v1.CreateWishlistItemResponse{
		ItemId: &common.ResourceID{
			Id: itemID,
		},
	}, nil
}

func (sa *serverAPI) DeleteWishlistItem(ctx context.Context, req *v1.DeleteWishlistItemRequest) (*v1.DeleteWishlistItemResponse, error) {
	err := sa.service.DeleteItem(ctx, req.ItemId.GetId())
	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Wishlist item not found")
		}
		return nil, status.Error(codes.Internal, "Failed to delete wishlist item")
	}

	return &v1.DeleteWishlistItemResponse{}, nil
}

func (sa *serverAPI) GetWishlistItem(ctx context.Context, req *v1.GetWishlistItemRequest) (*v1.GetWishlistItemResponse, error) {
	item, err := sa.service.GetItem(ctx, req.ItemId.GetId())
	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Wishlist item not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get wishlist item")
	}

	return &v1.GetWishlistItemResponse{
		WishlistItem: &v1.WishlistItem{
			ItemId:    item.ItemID,
			ProfileId: item.ProfileID,
			ProductId: item.ProductID,
			AddedAt:   timestamppb.New(item.AddedAt),
			IsActive:  item.IsActive,
		},
	}, nil
}

func (sa *serverAPI) ListWishlistItems(ctx context.Context, req *v1.ListWishlistItemsRequest) (*v1.ListWishlistItemsResponse, error) {
	items, err := sa.service.ListItems(ctx, req.GetProfileId())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to list wishlist items")
	}

	result := make([]*v1.WishlistItem, len(items))
	for i, item := range items {
		result[i] = &v1.WishlistItem{
			ItemId:    item.ItemID,
			ProfileId: item.ProfileID,
			ProductId: item.ProductID,
			AddedAt:   timestamppb.New(item.AddedAt),
			IsActive:  item.IsActive,
		}
	}

	return &v1.ListWishlistItemsResponse{
		WishlistItems: result,
	}, nil
}

func (sa *serverAPI) UpdateWishlistItem(ctx context.Context, req *v1.UpdateWishlistItemRequest) (*v1.UpdateWishlistItemResponse, error) {
	// Получаем текущий элемент
	currentItem, err := sa.service.GetItem(ctx, req.GetWishlistItem().GetItemId())
	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Wishlist item not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get wishlist item")
	}

	// Обновляем только is_active, остальные поля не меняются
	currentItem.IsActive = req.GetWishlistItem().GetIsActive()

	return &v1.UpdateWishlistItemResponse{}, nil
}
