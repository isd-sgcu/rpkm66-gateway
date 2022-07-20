// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: group.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GroupServiceClient is the client API for GroupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GroupServiceClient interface {
	FindOne(ctx context.Context, in *FindOneGroupRequest, opts ...grpc.CallOption) (*FindOneGroupResponse, error)
	FindByToken(ctx context.Context, in *FindByTokenGroupRequest, opts ...grpc.CallOption) (*FindByTokenGroupResponse, error)
	Update(ctx context.Context, in *UpdateGroupRequest, opts ...grpc.CallOption) (*UpdateGroupResponse, error)
	Join(ctx context.Context, in *JoinGroupRequest, opts ...grpc.CallOption) (*JoinGroupResponse, error)
	DeleteMember(ctx context.Context, in *DeleteMemberGroupRequest, opts ...grpc.CallOption) (*DeleteMemberGroupResponse, error)
	Leave(ctx context.Context, in *LeaveGroupRequest, opts ...grpc.CallOption) (*LeaveGroupResponse, error)
	SelectBaan(ctx context.Context, in *SelectBaanRequest, opts ...grpc.CallOption) (*SelectBaanResponse, error)
}

type groupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupServiceClient(cc grpc.ClientConnInterface) GroupServiceClient {
	return &groupServiceClient{cc}
}

func (c *groupServiceClient) FindOne(ctx context.Context, in *FindOneGroupRequest, opts ...grpc.CallOption) (*FindOneGroupResponse, error) {
	out := new(FindOneGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/FindOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) FindByToken(ctx context.Context, in *FindByTokenGroupRequest, opts ...grpc.CallOption) (*FindByTokenGroupResponse, error) {
	out := new(FindByTokenGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/FindByToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) Update(ctx context.Context, in *UpdateGroupRequest, opts ...grpc.CallOption) (*UpdateGroupResponse, error) {
	out := new(UpdateGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) Join(ctx context.Context, in *JoinGroupRequest, opts ...grpc.CallOption) (*JoinGroupResponse, error) {
	out := new(JoinGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) DeleteMember(ctx context.Context, in *DeleteMemberGroupRequest, opts ...grpc.CallOption) (*DeleteMemberGroupResponse, error) {
	out := new(DeleteMemberGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/DeleteMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) Leave(ctx context.Context, in *LeaveGroupRequest, opts ...grpc.CallOption) (*LeaveGroupResponse, error) {
	out := new(LeaveGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) SelectBaan(ctx context.Context, in *SelectBaanRequest, opts ...grpc.CallOption) (*SelectBaanResponse, error) {
	out := new(SelectBaanResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/SelectBaan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GroupServiceServer is the server API for GroupService service.
// All implementations must embed UnimplementedGroupServiceServer
// for forward compatibility
type GroupServiceServer interface {
	FindOne(context.Context, *FindOneGroupRequest) (*FindOneGroupResponse, error)
	FindByToken(context.Context, *FindByTokenGroupRequest) (*FindByTokenGroupResponse, error)
	Update(context.Context, *UpdateGroupRequest) (*UpdateGroupResponse, error)
	Join(context.Context, *JoinGroupRequest) (*JoinGroupResponse, error)
	DeleteMember(context.Context, *DeleteMemberGroupRequest) (*DeleteMemberGroupResponse, error)
	Leave(context.Context, *LeaveGroupRequest) (*LeaveGroupResponse, error)
	SelectBaan(context.Context, *SelectBaanRequest) (*SelectBaanResponse, error)
	mustEmbedUnimplementedGroupServiceServer()
}

// UnimplementedGroupServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGroupServiceServer struct {
}

func (UnimplementedGroupServiceServer) FindOne(context.Context, *FindOneGroupRequest) (*FindOneGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindOne not implemented")
}
func (UnimplementedGroupServiceServer) FindByToken(context.Context, *FindByTokenGroupRequest) (*FindByTokenGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByToken not implemented")
}
func (UnimplementedGroupServiceServer) Update(context.Context, *UpdateGroupRequest) (*UpdateGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedGroupServiceServer) Join(context.Context, *JoinGroupRequest) (*JoinGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedGroupServiceServer) DeleteMember(context.Context, *DeleteMemberGroupRequest) (*DeleteMemberGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMember not implemented")
}
func (UnimplementedGroupServiceServer) Leave(context.Context, *LeaveGroupRequest) (*LeaveGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedGroupServiceServer) SelectBaan(context.Context, *SelectBaanRequest) (*SelectBaanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SelectBaan not implemented")
}
func (UnimplementedGroupServiceServer) mustEmbedUnimplementedGroupServiceServer() {}

// UnsafeGroupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GroupServiceServer will
// result in compilation errors.
type UnsafeGroupServiceServer interface {
	mustEmbedUnimplementedGroupServiceServer()
}

func RegisterGroupServiceServer(s grpc.ServiceRegistrar, srv GroupServiceServer) {
	s.RegisterService(&GroupService_ServiceDesc, srv)
}

func _GroupService_FindOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindOneGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).FindOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/group.GroupService/FindOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).FindOne(ctx, req.(*FindOneGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_FindByToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByTokenGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).FindByToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/group.GroupService/FindByToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).FindByToken(ctx, req.(*FindByTokenGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/group.GroupService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).Update(ctx, req.(*UpdateGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/group.GroupService/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).Join(ctx, req.(*JoinGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_DeleteMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMemberGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).DeleteMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/group.GroupService/DeleteMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).DeleteMember(ctx, req.(*DeleteMemberGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/group.GroupService/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).Leave(ctx, req.(*LeaveGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_SelectBaan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SelectBaanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).SelectBaan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/group.GroupService/SelectBaan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).SelectBaan(ctx, req.(*SelectBaanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GroupService_ServiceDesc is the grpc.ServiceDesc for GroupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GroupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "group.GroupService",
	HandlerType: (*GroupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindOne",
			Handler:    _GroupService_FindOne_Handler,
		},
		{
			MethodName: "FindByToken",
			Handler:    _GroupService_FindByToken_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _GroupService_Update_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _GroupService_Join_Handler,
		},
		{
			MethodName: "DeleteMember",
			Handler:    _GroupService_DeleteMember_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _GroupService_Leave_Handler,
		},
		{
			MethodName: "SelectBaan",
			Handler:    _GroupService_SelectBaan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "group.proto",
}
