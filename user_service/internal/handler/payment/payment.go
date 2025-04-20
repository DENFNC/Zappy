package payment

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
)

type Payment interface {
	Create(ctx context.Context, profileID uint32, paymentToken string, isDefault bool) (uint32, error)
	GetByID(ctx context.Context, paymentID uint32) (*models.Payment, error)
	Update(ctx context.Context, paymentID uint32, profileID uint32, paymentToken string, isDefault bool) (uint32, error)
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
	panic("implement me!")
}

func (sa *serverAPI) GetPayment(ctx context.Context, req *v1.ResourceByIDRequest) (*v1.Payment, error) {
	panic("implement me!")
}

func (sa *serverAPI) UpdatePayment(ctx context.Context, req *v1.Payment) (*v1.ResourceID, error) {
	panic("implement me!")
}

func (sa *serverAPI) DeletePayment(ctx context.Context, req *v1.ResourceByIDRequest) (*v1.ResourceID, error) {
	panic("implement me!")
}

func (sa *serverAPI) ListPayments(ctx context.Context, req *v1.ListByProfileRequest) (*v1.ListPaymentResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) SetDefaultPayment(ctx context.Context, req *v1.SetDefaultPaymentRequest) (*v1.ResourceID, error) {
	panic("implement me!")
}
