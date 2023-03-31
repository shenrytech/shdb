// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.1
// source: shdb.proto

package shdb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ObjectServiceClient is the client API for ObjectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObjectServiceClient interface {
	List(ctx context.Context, in *ListReq, opts ...grpc.CallOption) (*ListRsp, error)
	Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*Object, error)
	Create(ctx context.Context, in *CreateReq, opts ...grpc.CallOption) (*Object, error)
	Update(ctx context.Context, in *UpdateReq, opts ...grpc.CallOption) (*Object, error)
	Delete(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*Object, error)
	GetSchema(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*descriptorpb.FileDescriptorSet, error)
	GetTypeNames(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetTypeNamesRsp, error)
}

type objectServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewObjectServiceClient(cc grpc.ClientConnInterface) ObjectServiceClient {
	return &objectServiceClient{cc}
}

func (c *objectServiceClient) List(ctx context.Context, in *ListReq, opts ...grpc.CallOption) (*ListRsp, error) {
	out := new(ListRsp)
	err := c.cc.Invoke(ctx, "/shdb.ObjectService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectServiceClient) Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*Object, error) {
	out := new(Object)
	err := c.cc.Invoke(ctx, "/shdb.ObjectService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectServiceClient) Create(ctx context.Context, in *CreateReq, opts ...grpc.CallOption) (*Object, error) {
	out := new(Object)
	err := c.cc.Invoke(ctx, "/shdb.ObjectService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectServiceClient) Update(ctx context.Context, in *UpdateReq, opts ...grpc.CallOption) (*Object, error) {
	out := new(Object)
	err := c.cc.Invoke(ctx, "/shdb.ObjectService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectServiceClient) Delete(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*Object, error) {
	out := new(Object)
	err := c.cc.Invoke(ctx, "/shdb.ObjectService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectServiceClient) GetSchema(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*descriptorpb.FileDescriptorSet, error) {
	out := new(descriptorpb.FileDescriptorSet)
	err := c.cc.Invoke(ctx, "/shdb.ObjectService/GetSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectServiceClient) GetTypeNames(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetTypeNamesRsp, error) {
	out := new(GetTypeNamesRsp)
	err := c.cc.Invoke(ctx, "/shdb.ObjectService/GetTypeNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ObjectServiceServer is the server API for ObjectService service.
// All implementations must embed UnimplementedObjectServiceServer
// for forward compatibility
type ObjectServiceServer interface {
	List(context.Context, *ListReq) (*ListRsp, error)
	Get(context.Context, *GetReq) (*Object, error)
	Create(context.Context, *CreateReq) (*Object, error)
	Update(context.Context, *UpdateReq) (*Object, error)
	Delete(context.Context, *DeleteReq) (*Object, error)
	GetSchema(context.Context, *emptypb.Empty) (*descriptorpb.FileDescriptorSet, error)
	GetTypeNames(context.Context, *emptypb.Empty) (*GetTypeNamesRsp, error)
	mustEmbedUnimplementedObjectServiceServer()
}

// UnimplementedObjectServiceServer must be embedded to have forward compatible implementations.
type UnimplementedObjectServiceServer struct {
}

func (UnimplementedObjectServiceServer) List(context.Context, *ListReq) (*ListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedObjectServiceServer) Get(context.Context, *GetReq) (*Object, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedObjectServiceServer) Create(context.Context, *CreateReq) (*Object, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedObjectServiceServer) Update(context.Context, *UpdateReq) (*Object, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedObjectServiceServer) Delete(context.Context, *DeleteReq) (*Object, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedObjectServiceServer) GetSchema(context.Context, *emptypb.Empty) (*descriptorpb.FileDescriptorSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchema not implemented")
}
func (UnimplementedObjectServiceServer) GetTypeNames(context.Context, *emptypb.Empty) (*GetTypeNamesRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTypeNames not implemented")
}
func (UnimplementedObjectServiceServer) mustEmbedUnimplementedObjectServiceServer() {}

// UnsafeObjectServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ObjectServiceServer will
// result in compilation errors.
type UnsafeObjectServiceServer interface {
	mustEmbedUnimplementedObjectServiceServer()
}

func RegisterObjectServiceServer(s grpc.ServiceRegistrar, srv ObjectServiceServer) {
	s.RegisterService(&ObjectService_ServiceDesc, srv)
}

func _ObjectService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shdb.ObjectService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServiceServer).List(ctx, req.(*ListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shdb.ObjectService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServiceServer).Get(ctx, req.(*GetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shdb.ObjectService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServiceServer).Create(ctx, req.(*CreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shdb.ObjectService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServiceServer).Update(ctx, req.(*UpdateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shdb.ObjectService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServiceServer).Delete(ctx, req.(*DeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectService_GetSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServiceServer).GetSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shdb.ObjectService/GetSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServiceServer).GetSchema(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectService_GetTypeNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServiceServer).GetTypeNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shdb.ObjectService/GetTypeNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServiceServer).GetTypeNames(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ObjectService_ServiceDesc is the grpc.ServiceDesc for ObjectService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ObjectService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shdb.ObjectService",
	HandlerType: (*ObjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _ObjectService_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ObjectService_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _ObjectService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ObjectService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ObjectService_Delete_Handler,
		},
		{
			MethodName: "GetSchema",
			Handler:    _ObjectService_GetSchema_Handler,
		},
		{
			MethodName: "GetTypeNames",
			Handler:    _ObjectService_GetTypeNames_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shdb.proto",
}
