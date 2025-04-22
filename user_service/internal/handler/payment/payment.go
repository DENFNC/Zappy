package payment

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Payment interface {
	Create(ctx context.Context, profileID uint32, paymentToken string) (uint32, error)
	GetByID(ctx context.Context, paymentID uint32) (*models.Payment, error)
	Update(ctx context.Context, paymentID uint32, profileID uint32, paymentToken string) (uint32, error)
	Delete(ctx context.Context, paymentID uint32) (uint32, error)
	List(ctx context.Context, profileID uint32) ([]*models.Payment, error)
	SetDefault(ctx context.Context, paymentID uint32, profileID uint32) (uint32, error)
}

type serverAPI struct {
	v1.UnimplementedPaymentServiceServer
	service Payment
}

func New(service Payment) *serverAPI {
	return &serverAPI{
		service: service,
	}
}

func (sa *serverAPI) Register(grpc *grpc.Server) {
	v1.RegisterPaymentServiceServer(grpc, sa)
}

func (sa *serverAPI) CreatePayment(ctx context.Context, req *v1.PaymentInput) (*v1.ResourceID, error) {
	payID, err := sa.service.Create(
		ctx,
		req.GetProfileId(),
		req.GetPaymentToken(),
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.ResourceID{
		Id: payID,
	}, nil
}

func (sa *serverAPI) GetPayment(ctx context.Context, req *v1.ResourceByIDRequest) (*v1.Payment, error) {
	payment, err := sa.service.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.Payment{
		PaymentId:    payment.PaymentID,
		ProfileId:    payment.ProfileID,
		PaymentToken: payment.PaymentToken,
		IsDefault:    payment.IsDefault,
	}, nil
}

func (sa *serverAPI) UpdatePayment(ctx context.Context, req *v1.Payment) (*v1.ResourceID, error) {
	panic("implement me!")
}

func (sa *serverAPI) DeletePayment(ctx context.Context, req *v1.ResourceByIDRequest) (*v1.ResourceID, error) {
	payID, err := sa.service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.ResourceID{
		Id: payID,
	}, nil
}

func (sa *serverAPI) ListPayments(ctx context.Context, req *v1.ListByProfileRequest) (*v1.ListPaymentResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) SetDefaultPayment(ctx context.Context, req *v1.SetDefaultPaymentRequest) (*v1.ResourceID, error) {
	panic("implement me!")
}
