package payment

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
	"github.com/DENFNC/Zappy/user_service/proto/gen/go/common/v1"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/go/payment/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Payment interface {
	CreatePayment(
		ctx context.Context,
		payment *models.Payment,
	) (string, error)
	PaymentGetByID(
		ctx context.Context,
		paymentID string,
	) (*models.Payment, error)
	UpdatePayment(
		ctx context.Context,
		uid string,
		payment *models.Payment,
	) (string, error)
	DeletePayment(
		ctx context.Context,
		paymentID string,
	) (string, error)
	ListPayments(
		ctx context.Context,
		profileID string,
		pageSize uint,
		pageToken string,
	) ([]*models.Payment, string, error)
	SetDefaultPayment(
		ctx context.Context,
		payment *models.Payment,
	) error
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

func (sa *serverAPI) GRPCRegister(grpc *grpc.Server) {
	v1.RegisterPaymentServiceServer(grpc, sa)
}

func (api *serverAPI) HTTPRegister(
	ctx context.Context,
	mux *runtime.ServeMux,
) {
	v1.RegisterPaymentServiceHandlerServer(ctx, mux, api)
}

func (sa *serverAPI) CreatePayment(ctx context.Context, req *v1.CreatePaymentRequest) (*v1.CreatePaymentResponse, error) {
	payID, err := sa.service.CreatePayment(
		ctx,
		&models.Payment{
			ProfileID:    req.GetProfileId(),
			PaymentToken: req.GetPaymentToken(),
			IsDefault:    req.GetIsDefault(),
		},
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.CreatePaymentResponse{
		Id: &common.ResourceID{
			Id: payID,
		},
	}, nil
}

func (sa *serverAPI) GetPayment(ctx context.Context, req *v1.GetPaymentRequest) (*v1.GetPaymentResponse, error) {
	payment, err := sa.service.PaymentGetByID(ctx, req.Id.GetId())
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.GetPaymentResponse{
		Payment: &v1.Payment{
			PaymentId:    payment.PaymentID,
			ProfileId:    payment.ProfileID,
			PaymentToken: payment.PaymentToken,
			IsDefault:    payment.IsDefault,
		},
	}, nil
}

func (sa *serverAPI) UpdatePayment(ctx context.Context, req *v1.UpdatePaymentRequest) (*v1.UpdatePaymentResponse, error) {
	payID, err := sa.service.UpdatePayment(
		ctx,
		req.PaymentId.GetId(),
		&models.Payment{
			PaymentToken: req.Payment.GetPaymentToken(),
		},
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

	return &v1.UpdatePaymentResponse{
		Id: &common.ResourceID{
			Id: payID,
		},
	}, nil
}

func (sa *serverAPI) DeletePayment(ctx context.Context, req *v1.DeletePaymentRequest) (*v1.DeletePaymentResponse, error) {
	payID, err := sa.service.DeletePayment(ctx, req.Id.GetId())
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.DeletePaymentResponse{
		Id: &common.ResourceID{
			Id: payID,
		},
	}, nil
}

func (sa *serverAPI) ListPayments(ctx context.Context, req *v1.ListPaymentsRequest) (*v1.ListPaymentsResponse, error) {
	afterPage := true
	items, pageToken, err := sa.service.ListPayments(
		ctx,
		req.GetProfileId(),
		uint(req.Pagination.GetPageSize()),
		req.Pagination.GetPageToken(),
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}
	if pageToken == "" {
		afterPage = false
	}

	v1Payments := make([]*v1.Payment, len(items))
	for i, item := range items {
		v1Payments[i] = &v1.Payment{
			PaymentId:    item.PaymentID,
			ProfileId:    item.ProfileID,
			PaymentToken: item.PaymentToken,
			IsDefault:    item.IsDefault,
		}
	}

	return &v1.ListPaymentsResponse{
		Payments: v1Payments,
		Pagination: &common.PaginationResponse{
			PageSize:  uint32(len(items)),
			PageToken: pageToken,
			AfterPage: afterPage,
		},
	}, nil
}

func (sa *serverAPI) SetDefaultPayment(ctx context.Context, req *v1.SetDefaultPaymentRequest) (*v1.SetDefaultPaymentResponse, error) {
	err := sa.service.SetDefaultPayment(
		ctx,
		&models.Payment{
			PaymentID: req.PaymentId.GetId(),
			ProfileID: req.GetProfileId(),
		},
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errpkg.ErrInternal.Message,
		)
	}

	return &v1.SetDefaultPaymentResponse{}, nil
}
