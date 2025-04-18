// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: payment.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PaymentMethodService_CreatePaymentMethod_FullMethodName = "/user.v1.PaymentMethodService/CreatePaymentMethod"
	PaymentMethodService_GetPaymentMethod_FullMethodName    = "/user.v1.PaymentMethodService/GetPaymentMethod"
	PaymentMethodService_UpdatePaymentMethod_FullMethodName = "/user.v1.PaymentMethodService/UpdatePaymentMethod"
	PaymentMethodService_DeletePaymentMethod_FullMethodName = "/user.v1.PaymentMethodService/DeletePaymentMethod"
	PaymentMethodService_ListPaymentMethods_FullMethodName  = "/user.v1.PaymentMethodService/ListPaymentMethods"
)

// PaymentMethodServiceClient is the client API for PaymentMethodService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentMethodServiceClient interface {
	CreatePaymentMethod(ctx context.Context, in *CreatePaymentMethodRequest, opts ...grpc.CallOption) (*PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, in *GetPaymentMethodRequest, opts ...grpc.CallOption) (*PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, in *UpdatePaymentMethodRequest, opts ...grpc.CallOption) (*PaymentMethod, error)
	DeletePaymentMethod(ctx context.Context, in *DeletePaymentMethodRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListPaymentMethods(ctx context.Context, in *ListPaymentMethodsRequest, opts ...grpc.CallOption) (*ListPaymentMethodsResponse, error)
}

type paymentMethodServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentMethodServiceClient(cc grpc.ClientConnInterface) PaymentMethodServiceClient {
	return &paymentMethodServiceClient{cc}
}

func (c *paymentMethodServiceClient) CreatePaymentMethod(ctx context.Context, in *CreatePaymentMethodRequest, opts ...grpc.CallOption) (*PaymentMethod, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PaymentMethod)
	err := c.cc.Invoke(ctx, PaymentMethodService_CreatePaymentMethod_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentMethodServiceClient) GetPaymentMethod(ctx context.Context, in *GetPaymentMethodRequest, opts ...grpc.CallOption) (*PaymentMethod, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PaymentMethod)
	err := c.cc.Invoke(ctx, PaymentMethodService_GetPaymentMethod_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentMethodServiceClient) UpdatePaymentMethod(ctx context.Context, in *UpdatePaymentMethodRequest, opts ...grpc.CallOption) (*PaymentMethod, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PaymentMethod)
	err := c.cc.Invoke(ctx, PaymentMethodService_UpdatePaymentMethod_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentMethodServiceClient) DeletePaymentMethod(ctx context.Context, in *DeletePaymentMethodRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PaymentMethodService_DeletePaymentMethod_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentMethodServiceClient) ListPaymentMethods(ctx context.Context, in *ListPaymentMethodsRequest, opts ...grpc.CallOption) (*ListPaymentMethodsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPaymentMethodsResponse)
	err := c.cc.Invoke(ctx, PaymentMethodService_ListPaymentMethods_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentMethodServiceServer is the server API for PaymentMethodService service.
// All implementations must embed UnimplementedPaymentMethodServiceServer
// for forward compatibility.
type PaymentMethodServiceServer interface {
	CreatePaymentMethod(context.Context, *CreatePaymentMethodRequest) (*PaymentMethod, error)
	GetPaymentMethod(context.Context, *GetPaymentMethodRequest) (*PaymentMethod, error)
	UpdatePaymentMethod(context.Context, *UpdatePaymentMethodRequest) (*PaymentMethod, error)
	DeletePaymentMethod(context.Context, *DeletePaymentMethodRequest) (*emptypb.Empty, error)
	ListPaymentMethods(context.Context, *ListPaymentMethodsRequest) (*ListPaymentMethodsResponse, error)
	mustEmbedUnimplementedPaymentMethodServiceServer()
}

// UnimplementedPaymentMethodServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPaymentMethodServiceServer struct{}

func (UnimplementedPaymentMethodServiceServer) CreatePaymentMethod(context.Context, *CreatePaymentMethodRequest) (*PaymentMethod, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePaymentMethod not implemented")
}
func (UnimplementedPaymentMethodServiceServer) GetPaymentMethod(context.Context, *GetPaymentMethodRequest) (*PaymentMethod, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaymentMethod not implemented")
}
func (UnimplementedPaymentMethodServiceServer) UpdatePaymentMethod(context.Context, *UpdatePaymentMethodRequest) (*PaymentMethod, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePaymentMethod not implemented")
}
func (UnimplementedPaymentMethodServiceServer) DeletePaymentMethod(context.Context, *DeletePaymentMethodRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePaymentMethod not implemented")
}
func (UnimplementedPaymentMethodServiceServer) ListPaymentMethods(context.Context, *ListPaymentMethodsRequest) (*ListPaymentMethodsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPaymentMethods not implemented")
}
func (UnimplementedPaymentMethodServiceServer) mustEmbedUnimplementedPaymentMethodServiceServer() {}
func (UnimplementedPaymentMethodServiceServer) testEmbeddedByValue()                              {}

// UnsafePaymentMethodServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentMethodServiceServer will
// result in compilation errors.
type UnsafePaymentMethodServiceServer interface {
	mustEmbedUnimplementedPaymentMethodServiceServer()
}

func RegisterPaymentMethodServiceServer(s grpc.ServiceRegistrar, srv PaymentMethodServiceServer) {
	// If the following call pancis, it indicates UnimplementedPaymentMethodServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PaymentMethodService_ServiceDesc, srv)
}

func _PaymentMethodService_CreatePaymentMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePaymentMethodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentMethodServiceServer).CreatePaymentMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentMethodService_CreatePaymentMethod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentMethodServiceServer).CreatePaymentMethod(ctx, req.(*CreatePaymentMethodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentMethodService_GetPaymentMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPaymentMethodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentMethodServiceServer).GetPaymentMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentMethodService_GetPaymentMethod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentMethodServiceServer).GetPaymentMethod(ctx, req.(*GetPaymentMethodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentMethodService_UpdatePaymentMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePaymentMethodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentMethodServiceServer).UpdatePaymentMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentMethodService_UpdatePaymentMethod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentMethodServiceServer).UpdatePaymentMethod(ctx, req.(*UpdatePaymentMethodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentMethodService_DeletePaymentMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePaymentMethodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentMethodServiceServer).DeletePaymentMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentMethodService_DeletePaymentMethod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentMethodServiceServer).DeletePaymentMethod(ctx, req.(*DeletePaymentMethodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentMethodService_ListPaymentMethods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPaymentMethodsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentMethodServiceServer).ListPaymentMethods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentMethodService_ListPaymentMethods_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentMethodServiceServer).ListPaymentMethods(ctx, req.(*ListPaymentMethodsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PaymentMethodService_ServiceDesc is the grpc.ServiceDesc for PaymentMethodService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PaymentMethodService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.PaymentMethodService",
	HandlerType: (*PaymentMethodServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePaymentMethod",
			Handler:    _PaymentMethodService_CreatePaymentMethod_Handler,
		},
		{
			MethodName: "GetPaymentMethod",
			Handler:    _PaymentMethodService_GetPaymentMethod_Handler,
		},
		{
			MethodName: "UpdatePaymentMethod",
			Handler:    _PaymentMethodService_UpdatePaymentMethod_Handler,
		},
		{
			MethodName: "DeletePaymentMethod",
			Handler:    _PaymentMethodService_DeletePaymentMethod_Handler,
		},
		{
			MethodName: "ListPaymentMethods",
			Handler:    _PaymentMethodService_ListPaymentMethods_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payment.proto",
}
