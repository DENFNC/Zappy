// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: shipping.proto

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
	ShippingAddressService_CreateShippingAddress_FullMethodName = "/user.v1.ShippingAddressService/CreateShippingAddress"
	ShippingAddressService_GetShippingAddress_FullMethodName    = "/user.v1.ShippingAddressService/GetShippingAddress"
	ShippingAddressService_UpdateShippingAddress_FullMethodName = "/user.v1.ShippingAddressService/UpdateShippingAddress"
	ShippingAddressService_DeleteShippingAddress_FullMethodName = "/user.v1.ShippingAddressService/DeleteShippingAddress"
	ShippingAddressService_ListShippingAddresses_FullMethodName = "/user.v1.ShippingAddressService/ListShippingAddresses"
)

// ShippingAddressServiceClient is the client API for ShippingAddressService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShippingAddressServiceClient interface {
	CreateShippingAddress(ctx context.Context, in *CreateShippingAddressRequest, opts ...grpc.CallOption) (*ShippingAddress, error)
	GetShippingAddress(ctx context.Context, in *GetShippingAddressRequest, opts ...grpc.CallOption) (*ShippingAddress, error)
	UpdateShippingAddress(ctx context.Context, in *UpdateShippingAddressRequest, opts ...grpc.CallOption) (*ShippingAddress, error)
	DeleteShippingAddress(ctx context.Context, in *DeleteShippingAddressRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListShippingAddresses(ctx context.Context, in *ListShippingAddressesRequest, opts ...grpc.CallOption) (*ListShippingAddressesResponse, error)
}

type shippingAddressServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShippingAddressServiceClient(cc grpc.ClientConnInterface) ShippingAddressServiceClient {
	return &shippingAddressServiceClient{cc}
}

func (c *shippingAddressServiceClient) CreateShippingAddress(ctx context.Context, in *CreateShippingAddressRequest, opts ...grpc.CallOption) (*ShippingAddress, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShippingAddress)
	err := c.cc.Invoke(ctx, ShippingAddressService_CreateShippingAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shippingAddressServiceClient) GetShippingAddress(ctx context.Context, in *GetShippingAddressRequest, opts ...grpc.CallOption) (*ShippingAddress, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShippingAddress)
	err := c.cc.Invoke(ctx, ShippingAddressService_GetShippingAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shippingAddressServiceClient) UpdateShippingAddress(ctx context.Context, in *UpdateShippingAddressRequest, opts ...grpc.CallOption) (*ShippingAddress, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShippingAddress)
	err := c.cc.Invoke(ctx, ShippingAddressService_UpdateShippingAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shippingAddressServiceClient) DeleteShippingAddress(ctx context.Context, in *DeleteShippingAddressRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ShippingAddressService_DeleteShippingAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shippingAddressServiceClient) ListShippingAddresses(ctx context.Context, in *ListShippingAddressesRequest, opts ...grpc.CallOption) (*ListShippingAddressesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListShippingAddressesResponse)
	err := c.cc.Invoke(ctx, ShippingAddressService_ListShippingAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShippingAddressServiceServer is the server API for ShippingAddressService service.
// All implementations must embed UnimplementedShippingAddressServiceServer
// for forward compatibility.
type ShippingAddressServiceServer interface {
	CreateShippingAddress(context.Context, *CreateShippingAddressRequest) (*ShippingAddress, error)
	GetShippingAddress(context.Context, *GetShippingAddressRequest) (*ShippingAddress, error)
	UpdateShippingAddress(context.Context, *UpdateShippingAddressRequest) (*ShippingAddress, error)
	DeleteShippingAddress(context.Context, *DeleteShippingAddressRequest) (*emptypb.Empty, error)
	ListShippingAddresses(context.Context, *ListShippingAddressesRequest) (*ListShippingAddressesResponse, error)
	mustEmbedUnimplementedShippingAddressServiceServer()
}

// UnimplementedShippingAddressServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedShippingAddressServiceServer struct{}

func (UnimplementedShippingAddressServiceServer) CreateShippingAddress(context.Context, *CreateShippingAddressRequest) (*ShippingAddress, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShippingAddress not implemented")
}
func (UnimplementedShippingAddressServiceServer) GetShippingAddress(context.Context, *GetShippingAddressRequest) (*ShippingAddress, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShippingAddress not implemented")
}
func (UnimplementedShippingAddressServiceServer) UpdateShippingAddress(context.Context, *UpdateShippingAddressRequest) (*ShippingAddress, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateShippingAddress not implemented")
}
func (UnimplementedShippingAddressServiceServer) DeleteShippingAddress(context.Context, *DeleteShippingAddressRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteShippingAddress not implemented")
}
func (UnimplementedShippingAddressServiceServer) ListShippingAddresses(context.Context, *ListShippingAddressesRequest) (*ListShippingAddressesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListShippingAddresses not implemented")
}
func (UnimplementedShippingAddressServiceServer) mustEmbedUnimplementedShippingAddressServiceServer() {
}
func (UnimplementedShippingAddressServiceServer) testEmbeddedByValue() {}

// UnsafeShippingAddressServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShippingAddressServiceServer will
// result in compilation errors.
type UnsafeShippingAddressServiceServer interface {
	mustEmbedUnimplementedShippingAddressServiceServer()
}

func RegisterShippingAddressServiceServer(s grpc.ServiceRegistrar, srv ShippingAddressServiceServer) {
	// If the following call pancis, it indicates UnimplementedShippingAddressServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ShippingAddressService_ServiceDesc, srv)
}

func _ShippingAddressService_CreateShippingAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShippingAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShippingAddressServiceServer).CreateShippingAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShippingAddressService_CreateShippingAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShippingAddressServiceServer).CreateShippingAddress(ctx, req.(*CreateShippingAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShippingAddressService_GetShippingAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShippingAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShippingAddressServiceServer).GetShippingAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShippingAddressService_GetShippingAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShippingAddressServiceServer).GetShippingAddress(ctx, req.(*GetShippingAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShippingAddressService_UpdateShippingAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateShippingAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShippingAddressServiceServer).UpdateShippingAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShippingAddressService_UpdateShippingAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShippingAddressServiceServer).UpdateShippingAddress(ctx, req.(*UpdateShippingAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShippingAddressService_DeleteShippingAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteShippingAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShippingAddressServiceServer).DeleteShippingAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShippingAddressService_DeleteShippingAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShippingAddressServiceServer).DeleteShippingAddress(ctx, req.(*DeleteShippingAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShippingAddressService_ListShippingAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListShippingAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShippingAddressServiceServer).ListShippingAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShippingAddressService_ListShippingAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShippingAddressServiceServer).ListShippingAddresses(ctx, req.(*ListShippingAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShippingAddressService_ServiceDesc is the grpc.ServiceDesc for ShippingAddressService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShippingAddressService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.ShippingAddressService",
	HandlerType: (*ShippingAddressServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShippingAddress",
			Handler:    _ShippingAddressService_CreateShippingAddress_Handler,
		},
		{
			MethodName: "GetShippingAddress",
			Handler:    _ShippingAddressService_GetShippingAddress_Handler,
		},
		{
			MethodName: "UpdateShippingAddress",
			Handler:    _ShippingAddressService_UpdateShippingAddress_Handler,
		},
		{
			MethodName: "DeleteShippingAddress",
			Handler:    _ShippingAddressService_DeleteShippingAddress_Handler,
		},
		{
			MethodName: "ListShippingAddresses",
			Handler:    _ShippingAddressService_ListShippingAddresses_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shipping.proto",
}
