package payment

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

type Payment interface {
	Create(ctx context.Context, profileID string, paymentToken string) (string, error)
	GetByID(ctx context.Context, paymentID string) (*models.Payment, error)
	Update(ctx context.Context, paymentID string, profileID string, paymentToken string) (string, error)
	Delete(ctx context.Context, paymentID string) (string, error)
	List(ctx context.Context, profileID string) ([]models.Payment, error)
	SetDefault(ctx context.Context, paymentID string, profileID string) (string, error)
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

func (sa *serverAPI) UpdatePayment(ctx context.Context, req *v1.UpdatePaymentRequest) (*v1.ResourceID, error) {
	payID, err := sa.service.Update(
		ctx,
		req.GetPaymentId(),
		req.GetPayment().GetProfileId(),
		req.GetPayment().GetPaymentToken(),
	)
	if err != nil {
		switch {
		case errors.Is(err, errpkg.ErrNotFound):
			return nil, status.Error(
				codes.NotFound,
				errpkg.ErrNotFound.Message,
			)
		default:
			return nil, status.Error(
				codes.Internal,
				errpkg.ErrInternal.Message,
			)
		}
	}

	return &v1.ResourceID{
		Id: payID,
	}, nil
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
	payments, err := sa.service.List(ctx, req.GetProfileId())
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	v1Payments := make([]*v1.Payment, len(payments))
	for i, p := range payments {
		v1Payments[i] = &v1.Payment{
			PaymentId:    p.PaymentID,
			ProfileId:    p.ProfileID,
			PaymentToken: p.PaymentToken,
			IsDefault:    p.IsDefault,
		}
	}

	return &v1.ListPaymentResponse{
		Payments: v1Payments,
	}, nil
}

func (sa *serverAPI) SetDefaultPayment(ctx context.Context, req *v1.SetDefaultPaymentRequest) (*v1.ResourceID, error) {
	payID, err := sa.service.SetDefault(
		ctx,
		req.GetPaymentId(),
		req.GetProfileId(),
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
