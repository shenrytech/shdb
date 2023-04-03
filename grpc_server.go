// Copyright 2023 Shenry Tech AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shdb

import (
	"bytes"
	"context"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	UnimplementedObjectServiceServer
	ctx     context.Context
	gs      *grpc.Server
	typeReg *TypeRegistry
}

func NewServer(ctx context.Context, grpcServer *grpc.Server, typeReg *TypeRegistry) *Server {
	s := &Server{
		ctx:     ctx,
		gs:      grpcServer,
		typeReg: typeReg,
	}
	RegisterObjectServiceServer(grpcServer, s)
	return s
}

func (s *Server) List(ctx context.Context, req *ListReq) (*ListRsp, error) {
	list, nextPageToken, err := List[IObject](ctx, [4]byte(req.Type), req.PageSize, req.PageToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed listing objects")
	}

	kv, err := Marshal(list...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed listing objects")
	}
	rsp := &ListRsp{Objects: make([]*Object, 0), NextPageToken: nextPageToken}
	for _, v := range kv {
		rsp.Objects = append(rsp.Objects, &Object{Key: v.Key(), Value: v.Value})
	}
	return rsp, nil
}

func (s *Server) Get(ctx context.Context, req *GetReq) (*Object, error) {
	kv, err := get(*req.Ref.TypeId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed retrieve an object")
	}
	return &Object{Key: kv.Key(), Value: kv.Value}, nil

}

func (s *Server) Create(ctx context.Context, req *CreateReq) (*Object, error) {
	o, err := New[IObject]([4]byte(req.Type))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create object %v", err)
	}
	kv, err := Marshal(o)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create object %v", err)
	}
	return &Object{Key: kv[0].Key(), Value: kv[0].Value}, nil
}

func (s *Server) Update(ctx context.Context, req *UpdateReq) (*Object, error) {
	kv, err := get(*MarshalTypeId(req.Object.Key))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update object %v", err)
	}
	obj, err := Unmarshal[IObject](*kv)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update object %v", err)
	}
	ret, err := Update(obj.GetMetadata().TypeId(), func(obj IObject) (IObject, error) {
		return proto.Clone(obj).(IObject), nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update object %v", err)
	}
	kvs, err := Marshal(ret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update object %v", err)
	}
	return &Object{Key: kvs[0].Key(), Value: kvs[0].Value}, nil

}
func (s *Server) Delete(ctx context.Context, req *DeleteReq) (*Object, error) {
	obj, err := Delete[IObject](*req.Ref.TypeId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete object %v", err)
	}
	kvs, err := Marshal(obj)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete object %v", err)
	}
	return &Object{Key: kvs[0].Key(), Value: kvs[0].Value}, nil
}

func (s *Server) GetSchema(ctx context.Context, req *emptypb.Empty) (*descriptorpb.FileDescriptorSet, error) {
	return s.typeReg.GetFileDescriptorSet(), nil
}

func (s *Server) GetTypeNames(ctx context.Context, req *emptypb.Empty) (*GetTypeNamesRsp, error) {
	rsp := &GetTypeNamesRsp{TypeAliases: []*GetTypeNamesRsp_TypeAliases{}}
	tns := s.typeReg.GetTypeNames()
	for k, v := range tns {
		rsp.TypeAliases = append(rsp.TypeAliases, &GetTypeNamesRsp_TypeAliases{
			Fullname: k,
			Aliases:  v,
		})
	}
	return rsp, nil
}

func (s *Server) StreamRefs(req *StreamRefReq, stream ObjectService_StreamRefsServer) error {
	selector := func(obj *ObjRef) bool {
		if bytes.Equal(req.TypeKey, TypeKeyAll[:]) {
			return true
		}
		return bytes.Equal(req.TypeKey, obj.Type)
	}
	refs, _, err := SearchRef(stream.Context(), selector, 1000000, "")
	if err != nil {
		return status.Errorf(codes.Internal, "query ref failed")
	}
	for _, v := range refs {
		if err := stream.Send(v); err != nil {
			return err
		}
	}
	return nil
}
