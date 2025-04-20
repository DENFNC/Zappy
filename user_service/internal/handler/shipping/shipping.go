package shipping

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Shipping interface {
	Create(
		ctx context.Context,
		s *models.Shipping,
	) (uint32, error)
	GetByID(
		ctx context.Context,
		id uint32,
	) (*models.Shipping, error)
	ListByProfile(
		ctx context.Context,
		profileID uint32,
	) ([]models.Shipping, error)
	SetDefault(
		ctx context.Context,
		addressID, profileID uint32,
	) (uint32, error)
	Delete(
		ctx context.Context,
		id uint32,
	) (uint32, error)
}

type serverAPI struct {
	v1.UnimplementedShippingServiceServer
	service Shipping
}

func New(service Shipping) *serverAPI {
	return &serverAPI{
		service: service,
	}
}

func (sa *serverAPI) Register(grpc *grpc.Server) {
	v1.RegisterShippingServiceServer(grpc, sa)
}

func (sa *serverAPI) CreateShipping(ctx context.Context, req *v1.CreateShippingRequest) (*v1.CreateShippingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			errpkg.ErrInvalidArgument.Message,
		)
	}

	ship := models.NewShipping(
		req.GetAddress().ProfileId,
		req.GetAddress().Country,
		req.GetAddress().City,
		req.GetAddress().Street,
		req.GetAddress().PostalCode,
	)

	addrID, err := sa.service.Create(ctx, ship)
	if err != nil {
		switch {
		case errors.Is(err, errpkg.ErrConstraint):
			return nil, status.Error(
				codes.NotFound,
				errpkg.ErrConstraint.Message,
			)
		case errors.Is(err, errpkg.ErrUniqueViolation):
			return nil, status.Error(
				codes.AlreadyExists,
				errpkg.ErrUniqueViolation.Message,
			)
		default:
			return nil, status.Error(
				codes.Internal,
				errpkg.ErrInternal.Message,
			)
		}
	}

	return &v1.CreateShippingResponse{
		Id: &v1.ShippingId{
			XId: addrID,
		},
	}, nil
}

func (sa *serverAPI) DeleteShipping(ctx context.Context, req *v1.DeleteShippingRequest) (*v1.DeleteShippingResponse, error) {
	addrID, err := sa.service.Delete(
		ctx,
		req.GetId().XId,
	)
	if err != nil {
		switch {
		case errors.Is(err, errpkg.ErrNotFound):
			return nil, status.Error(
				codes.NotFound,
				errpkg.ErrNotFound.Message,
			)
		}
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.DeleteShippingResponse{
		Id: &v1.ShippingId{
			XId: addrID,
		},
	}, nil
}

func (sa *serverAPI) GetShipping(ctx context.Context, req *v1.GetShippingRequest) (*v1.GetShippingResponse, error) {
	ship, err := sa.service.GetByID(
		ctx,
		req.GetId().XId,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.GetShippingResponse{
		Address: &v1.Shipping{
			XId:        ship.AddressID,
			ProfileId:  ship.ProfileID,
			Country:    ship.Country,
			City:       ship.City,
			Street:     ship.Street,
			PostalCode: ship.PostalCode,
			IsDefault:  ship.IsDefault,
		},
	}, nil
}

func (sa *serverAPI) ListShipping(ctx context.Context, req *v1.ListShippingRequest) (*v1.ListShippingResponse, error) {
	list, err := sa.service.ListByProfile(ctx, req.GetProfileId())
	if err != nil {
		return nil, status.Error(codes.Internal, errpkg.ErrInternal.Message)
	}

	out := make([]*v1.Shipping, len(list))
	for i, addr := range list {
		out[i] = &v1.Shipping{
			XId:        addr.AddressID,
			ProfileId:  addr.ProfileID,
			Country:    addr.Country,
			City:       addr.City,
			Street:     addr.Street,
			PostalCode: addr.PostalCode,
			IsDefault:  addr.IsDefault,
		}
	}

	return &v1.ListShippingResponse{
		Es: out,
	}, nil
}

func (sa *serverAPI) UpdateShipping(context.Context, *v1.UpdateShippingRequest) (*v1.UpdateShippingResponse, error) {
	panic("implement me!")
}
