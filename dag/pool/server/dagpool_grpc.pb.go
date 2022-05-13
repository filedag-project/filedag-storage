// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: dagpool.proto

package server

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

// DagPoolClient is the client API for DagPool service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DagPoolClient interface {
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddReply, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
}

type dagPoolClient struct {
	cc grpc.ClientConnInterface
}

func NewDagPoolClient(cc grpc.ClientConnInterface) DagPoolClient {
	return &dagPoolClient{cc}
}

func (c *dagPoolClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddReply, error) {
	out := new(AddReply)
	err := c.cc.Invoke(ctx, "/proto.DagPool/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dagPoolClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, "/proto.DagPool/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dagPoolClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, "/proto.DagPool/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DagPoolServer is the server API for DagPool service.
// All implementations must embed UnimplementedDagPoolServer
// for forward compatibility
type DagPoolServer interface {
	Add(context.Context, *AddRequest) (*AddReply, error)
	Get(context.Context, *GetRequest) (*GetReply, error)
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
	mustEmbedUnimplementedDagPoolServer()
}

// UnimplementedDagPoolServer must be embedded to have forward compatible implementations.
type UnimplementedDagPoolServer struct {
}

func (UnimplementedDagPoolServer) Add(context.Context, *AddRequest) (*AddReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedDagPoolServer) Get(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedDagPoolServer) Delete(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDagPoolServer) mustEmbedUnimplementedDagPoolServer() {}

// UnsafeDagPoolServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DagPoolServer will
// result in compilation errors.
type UnsafeDagPoolServer interface {
	mustEmbedUnimplementedDagPoolServer()
}

func RegisterDagPoolServer(s grpc.ServiceRegistrar, srv DagPoolServer) {
	s.RegisterService(&DagPool_ServiceDesc, srv)
}

func _DagPool_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DagPoolServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DagPool/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DagPoolServer).Add(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DagPool_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DagPoolServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DagPool/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DagPoolServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DagPool_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DagPoolServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DagPool/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DagPoolServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DagPool_ServiceDesc is the grpc.ServiceDesc for DagPool service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DagPool_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DagPool",
	HandlerType: (*DagPoolServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _DagPool_Add_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DagPool_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DagPool_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dagpool.proto",
}