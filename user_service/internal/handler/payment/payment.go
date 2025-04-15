package payment

import (
	"context"

	v1 "github.com/DENFNC/Zappy/user_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Payment interface{}

type serverAPI struct {
	v1.UnimplementedPaymentMethodServiceServer
	service Payment
}

func New(service Payment) *serverAPI {
	return &serverAPI{
		service: service,
	}
}

func (sa *serverAPI) Register(grpc *grpc.Server) {
	v1.RegisterPaymentMethodServiceServer(grpc, sa)
}

func (sa *serverAPI) CreatePaymentMethod(context.Context, *v1.CreatePaymentMethodRequest) (*v1.PaymentMethod, error) {
	panic("implement me!")
}

func (sa *serverAPI) DeletePaymentMethod(context.Context, *v1.DeletePaymentMethodRequest) (*emptypb.Empty, error) {
	panic("implement me!")
}

func (sa *serverAPI) GetPaymentMethod(context.Context, *v1.GetPaymentMethodRequest) (*v1.PaymentMethod, error) {
	panic("implement me!")
}

func (sa *serverAPI) ListPaymentMethods(context.Context, *v1.ListPaymentMethodsRequest) (*v1.ListPaymentMethodsResponse, error) {
	panic("implement me!")
}

func (sa *serverAPI) UpdatePaymentMethod(context.Context, *v1.UpdatePaymentMethodRequest) (*v1.PaymentMethod, error) {
	panic("implement me!")
}
