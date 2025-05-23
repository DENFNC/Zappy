// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: profile.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UserProfileService_CreateProfile_FullMethodName = "/user.v1.UserProfileService/CreateProfile"
	UserProfileService_GetProfile_FullMethodName    = "/user.v1.UserProfileService/GetProfile"
	UserProfileService_UpdateProfile_FullMethodName = "/user.v1.UserProfileService/UpdateProfile"
	UserProfileService_DeleteProfile_FullMethodName = "/user.v1.UserProfileService/DeleteProfile"
	UserProfileService_ListProfiles_FullMethodName  = "/user.v1.UserProfileService/ListProfiles"
)

// UserProfileServiceClient is the client API for UserProfileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserProfileServiceClient interface {
	CreateProfile(ctx context.Context, in *CreateProfileRequest, opts ...grpc.CallOption) (*ProfileIDResponse, error)
	GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*Profile, error)
	UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*ProfileIDResponse, error)
	DeleteProfile(ctx context.Context, in *DeleteProfileRequest, opts ...grpc.CallOption) (*ProfileIDResponse, error)
	ListProfiles(ctx context.Context, in *ListProfilesRequest, opts ...grpc.CallOption) (*ListProfilesResponse, error)
}

type userProfileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserProfileServiceClient(cc grpc.ClientConnInterface) UserProfileServiceClient {
	return &userProfileServiceClient{cc}
}

func (c *userProfileServiceClient) CreateProfile(ctx context.Context, in *CreateProfileRequest, opts ...grpc.CallOption) (*ProfileIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProfileIDResponse)
	err := c.cc.Invoke(ctx, UserProfileService_CreateProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userProfileServiceClient) GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*Profile, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Profile)
	err := c.cc.Invoke(ctx, UserProfileService_GetProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userProfileServiceClient) UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*ProfileIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProfileIDResponse)
	err := c.cc.Invoke(ctx, UserProfileService_UpdateProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userProfileServiceClient) DeleteProfile(ctx context.Context, in *DeleteProfileRequest, opts ...grpc.CallOption) (*ProfileIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProfileIDResponse)
	err := c.cc.Invoke(ctx, UserProfileService_DeleteProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userProfileServiceClient) ListProfiles(ctx context.Context, in *ListProfilesRequest, opts ...grpc.CallOption) (*ListProfilesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListProfilesResponse)
	err := c.cc.Invoke(ctx, UserProfileService_ListProfiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserProfileServiceServer is the server API for UserProfileService service.
// All implementations must embed UnimplementedUserProfileServiceServer
// for forward compatibility.
type UserProfileServiceServer interface {
	CreateProfile(context.Context, *CreateProfileRequest) (*ProfileIDResponse, error)
	GetProfile(context.Context, *GetProfileRequest) (*Profile, error)
	UpdateProfile(context.Context, *UpdateProfileRequest) (*ProfileIDResponse, error)
	DeleteProfile(context.Context, *DeleteProfileRequest) (*ProfileIDResponse, error)
	ListProfiles(context.Context, *ListProfilesRequest) (*ListProfilesResponse, error)
	mustEmbedUnimplementedUserProfileServiceServer()
}

// UnimplementedUserProfileServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserProfileServiceServer struct{}

func (UnimplementedUserProfileServiceServer) CreateProfile(context.Context, *CreateProfileRequest) (*ProfileIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProfile not implemented")
}
func (UnimplementedUserProfileServiceServer) GetProfile(context.Context, *GetProfileRequest) (*Profile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}
func (UnimplementedUserProfileServiceServer) UpdateProfile(context.Context, *UpdateProfileRequest) (*ProfileIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}
func (UnimplementedUserProfileServiceServer) DeleteProfile(context.Context, *DeleteProfileRequest) (*ProfileIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProfile not implemented")
}
func (UnimplementedUserProfileServiceServer) ListProfiles(context.Context, *ListProfilesRequest) (*ListProfilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProfiles not implemented")
}
func (UnimplementedUserProfileServiceServer) mustEmbedUnimplementedUserProfileServiceServer() {}
func (UnimplementedUserProfileServiceServer) testEmbeddedByValue()                            {}

// UnsafeUserProfileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserProfileServiceServer will
// result in compilation errors.
type UnsafeUserProfileServiceServer interface {
	mustEmbedUnimplementedUserProfileServiceServer()
}

func RegisterUserProfileServiceServer(s grpc.ServiceRegistrar, srv UserProfileServiceServer) {
	// If the following call pancis, it indicates UnimplementedUserProfileServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserProfileService_ServiceDesc, srv)
}

func _UserProfileService_CreateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserProfileServiceServer).CreateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserProfileService_CreateProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserProfileServiceServer).CreateProfile(ctx, req.(*CreateProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserProfileService_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserProfileServiceServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserProfileService_GetProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserProfileServiceServer).GetProfile(ctx, req.(*GetProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserProfileService_UpdateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserProfileServiceServer).UpdateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserProfileService_UpdateProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserProfileServiceServer).UpdateProfile(ctx, req.(*UpdateProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserProfileService_DeleteProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserProfileServiceServer).DeleteProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserProfileService_DeleteProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserProfileServiceServer).DeleteProfile(ctx, req.(*DeleteProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserProfileService_ListProfiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProfilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserProfileServiceServer).ListProfiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserProfileService_ListProfiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserProfileServiceServer).ListProfiles(ctx, req.(*ListProfilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserProfileService_ServiceDesc is the grpc.ServiceDesc for UserProfileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.UserProfileService",
	HandlerType: (*UserProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProfile",
			Handler:    _UserProfileService_CreateProfile_Handler,
		},
		{
			MethodName: "GetProfile",
			Handler:    _UserProfileService_GetProfile_Handler,
		},
		{
			MethodName: "UpdateProfile",
			Handler:    _UserProfileService_UpdateProfile_Handler,
		},
		{
			MethodName: "DeleteProfile",
			Handler:    _UserProfileService_DeleteProfile_Handler,
		},
		{
			MethodName: "ListProfiles",
			Handler:    _UserProfileService_ListProfiles_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}
