// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: wishlist.proto

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
	WishlistItemService_CreateWishlistItem_FullMethodName = "/user.v1.WishlistItemService/CreateWishlistItem"
	WishlistItemService_GetWishlistItem_FullMethodName    = "/user.v1.WishlistItemService/GetWishlistItem"
	WishlistItemService_UpdateWishlistItem_FullMethodName = "/user.v1.WishlistItemService/UpdateWishlistItem"
	WishlistItemService_DeleteWishlistItem_FullMethodName = "/user.v1.WishlistItemService/DeleteWishlistItem"
	WishlistItemService_ListWishlistItems_FullMethodName  = "/user.v1.WishlistItemService/ListWishlistItems"
)

// WishlistItemServiceClient is the client API for WishlistItemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WishlistItemServiceClient interface {
	CreateWishlistItem(ctx context.Context, in *CreateWishlistItemRequest, opts ...grpc.CallOption) (*WishlistItem, error)
	GetWishlistItem(ctx context.Context, in *GetWishlistItemRequest, opts ...grpc.CallOption) (*WishlistItem, error)
	UpdateWishlistItem(ctx context.Context, in *UpdateWishlistItemRequest, opts ...grpc.CallOption) (*WishlistItem, error)
	DeleteWishlistItem(ctx context.Context, in *DeleteWishlistItemRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListWishlistItems(ctx context.Context, in *ListWishlistItemsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[WishlistItem], error)
}

type wishlistItemServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWishlistItemServiceClient(cc grpc.ClientConnInterface) WishlistItemServiceClient {
	return &wishlistItemServiceClient{cc}
}

func (c *wishlistItemServiceClient) CreateWishlistItem(ctx context.Context, in *CreateWishlistItemRequest, opts ...grpc.CallOption) (*WishlistItem, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WishlistItem)
	err := c.cc.Invoke(ctx, WishlistItemService_CreateWishlistItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wishlistItemServiceClient) GetWishlistItem(ctx context.Context, in *GetWishlistItemRequest, opts ...grpc.CallOption) (*WishlistItem, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WishlistItem)
	err := c.cc.Invoke(ctx, WishlistItemService_GetWishlistItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wishlistItemServiceClient) UpdateWishlistItem(ctx context.Context, in *UpdateWishlistItemRequest, opts ...grpc.CallOption) (*WishlistItem, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WishlistItem)
	err := c.cc.Invoke(ctx, WishlistItemService_UpdateWishlistItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wishlistItemServiceClient) DeleteWishlistItem(ctx context.Context, in *DeleteWishlistItemRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, WishlistItemService_DeleteWishlistItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wishlistItemServiceClient) ListWishlistItems(ctx context.Context, in *ListWishlistItemsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[WishlistItem], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &WishlistItemService_ServiceDesc.Streams[0], WishlistItemService_ListWishlistItems_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ListWishlistItemsRequest, WishlistItem]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type WishlistItemService_ListWishlistItemsClient = grpc.ServerStreamingClient[WishlistItem]

// WishlistItemServiceServer is the server API for WishlistItemService service.
// All implementations must embed UnimplementedWishlistItemServiceServer
// for forward compatibility.
type WishlistItemServiceServer interface {
	CreateWishlistItem(context.Context, *CreateWishlistItemRequest) (*WishlistItem, error)
	GetWishlistItem(context.Context, *GetWishlistItemRequest) (*WishlistItem, error)
	UpdateWishlistItem(context.Context, *UpdateWishlistItemRequest) (*WishlistItem, error)
	DeleteWishlistItem(context.Context, *DeleteWishlistItemRequest) (*emptypb.Empty, error)
	ListWishlistItems(*ListWishlistItemsRequest, grpc.ServerStreamingServer[WishlistItem]) error
	mustEmbedUnimplementedWishlistItemServiceServer()
}

// UnimplementedWishlistItemServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedWishlistItemServiceServer struct{}

func (UnimplementedWishlistItemServiceServer) CreateWishlistItem(context.Context, *CreateWishlistItemRequest) (*WishlistItem, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWishlistItem not implemented")
}
func (UnimplementedWishlistItemServiceServer) GetWishlistItem(context.Context, *GetWishlistItemRequest) (*WishlistItem, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWishlistItem not implemented")
}
func (UnimplementedWishlistItemServiceServer) UpdateWishlistItem(context.Context, *UpdateWishlistItemRequest) (*WishlistItem, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateWishlistItem not implemented")
}
func (UnimplementedWishlistItemServiceServer) DeleteWishlistItem(context.Context, *DeleteWishlistItemRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteWishlistItem not implemented")
}
func (UnimplementedWishlistItemServiceServer) ListWishlistItems(*ListWishlistItemsRequest, grpc.ServerStreamingServer[WishlistItem]) error {
	return status.Errorf(codes.Unimplemented, "method ListWishlistItems not implemented")
}
func (UnimplementedWishlistItemServiceServer) mustEmbedUnimplementedWishlistItemServiceServer() {}
func (UnimplementedWishlistItemServiceServer) testEmbeddedByValue()                             {}

// UnsafeWishlistItemServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WishlistItemServiceServer will
// result in compilation errors.
type UnsafeWishlistItemServiceServer interface {
	mustEmbedUnimplementedWishlistItemServiceServer()
}

func RegisterWishlistItemServiceServer(s grpc.ServiceRegistrar, srv WishlistItemServiceServer) {
	// If the following call pancis, it indicates UnimplementedWishlistItemServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&WishlistItemService_ServiceDesc, srv)
}

func _WishlistItemService_CreateWishlistItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWishlistItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WishlistItemServiceServer).CreateWishlistItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WishlistItemService_CreateWishlistItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WishlistItemServiceServer).CreateWishlistItem(ctx, req.(*CreateWishlistItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WishlistItemService_GetWishlistItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWishlistItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WishlistItemServiceServer).GetWishlistItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WishlistItemService_GetWishlistItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WishlistItemServiceServer).GetWishlistItem(ctx, req.(*GetWishlistItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WishlistItemService_UpdateWishlistItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateWishlistItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WishlistItemServiceServer).UpdateWishlistItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WishlistItemService_UpdateWishlistItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WishlistItemServiceServer).UpdateWishlistItem(ctx, req.(*UpdateWishlistItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WishlistItemService_DeleteWishlistItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteWishlistItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WishlistItemServiceServer).DeleteWishlistItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WishlistItemService_DeleteWishlistItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WishlistItemServiceServer).DeleteWishlistItem(ctx, req.(*DeleteWishlistItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WishlistItemService_ListWishlistItems_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListWishlistItemsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WishlistItemServiceServer).ListWishlistItems(m, &grpc.GenericServerStream[ListWishlistItemsRequest, WishlistItem]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type WishlistItemService_ListWishlistItemsServer = grpc.ServerStreamingServer[WishlistItem]

// WishlistItemService_ServiceDesc is the grpc.ServiceDesc for WishlistItemService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WishlistItemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.WishlistItemService",
	HandlerType: (*WishlistItemServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateWishlistItem",
			Handler:    _WishlistItemService_CreateWishlistItem_Handler,
		},
		{
			MethodName: "GetWishlistItem",
			Handler:    _WishlistItemService_GetWishlistItem_Handler,
		},
		{
			MethodName: "UpdateWishlistItem",
			Handler:    _WishlistItemService_UpdateWishlistItem_Handler,
		},
		{
			MethodName: "DeleteWishlistItem",
			Handler:    _WishlistItemService_DeleteWishlistItem_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListWishlistItems",
			Handler:       _WishlistItemService_ListWishlistItems_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "wishlist.proto",
}
