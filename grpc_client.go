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
	"context"
	"io"
	"log"

	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	cc      *grpc.ClientConn
	cli     ObjectServiceClient
	ctx     context.Context
	typeReg *TypeRegistry
}

// NewClient returns a new client for use with the API
func NewClient(ctx context.Context, cc *grpc.ClientConn) *Client {
	return &Client{ctx: ctx, cc: cc, cli: NewObjectServiceClient(cc), typeReg: nil}
}

func (c *Client) Get(ctx context.Context, tid TypeId) (IObject, error) {
	ref, err := UnmarshalObjRef(tid.Key())
	if err != nil {
		return nil, err
	}
	o, err := c.cli.Get(ctx, &GetReq{Ref: ref})
	if err != nil {
		return nil, err
	}
	kv := KeyVal{TypeId: *ref.TypeId(), Value: o.Value}
	return Unmarshal[IObject](kv)
}

func (c *Client) List(ctx context.Context, tid TypeId) ([]IObject, error) {
	ref, err := UnmarshalObjRef(tid.Key())
	if err != nil {
		return nil, err
	}
	rsp, err := c.cli.List(ctx, &ListReq{PageSize: 100000, PageToken: ""})
	if err != nil {
		return nil, err
	}
	res := []IObject{}
	for _, v := range rsp.Objects {
		kv := KeyVal{TypeId: *ref.TypeId(), Value: v.Value}
		obj, err := Unmarshal[IObject](kv)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}
	return res, nil
}

func (c *Client) Delete(ctx context.Context, tid TypeId) (IObject, error) {
	ref, err := UnmarshalObjRef(tid.Key())
	if err != nil {
		return nil, err
	}
	rsp, err := c.cli.Delete(ctx, &DeleteReq{Ref: ref})
	if err != nil {
		return nil, err
	}
	kv := KeyVal{TypeId: *ref.TypeId(), Value: rsp.Value}
	return Unmarshal[IObject](kv)
}

func (c *Client) Create(ctx context.Context, typ TypeKey) (IObject, error) {
	rsp, err := c.cli.Create(ctx, &CreateReq{Type: typ[:]})
	if err != nil {
		return nil, err
	}
	kv := KeyVal{TypeId: *MarshalTypeId(rsp.Key), Value: rsp.Value}
	return Unmarshal[IObject](kv)
}

func (c *Client) Update(ctx context.Context, obj IObject) (IObject, error) {
	kvs, err := Marshal(obj)
	if err != nil {
		return nil, err
	}
	o := &Object{Key: kvs[0].Key(), Value: kvs[0].Value}
	rsp, err := c.cli.Update(ctx, &UpdateReq{Object: o})
	if err != nil {
		return nil, err
	}
	kv := KeyVal{TypeId: *MarshalTypeId(rsp.Key), Value: o.Value}
	return Unmarshal[IObject](kv)
}

func (c *Client) TypeRegistry() *TypeRegistry {
	if c.typeReg == nil {
		schema, err := c.cli.GetSchema(c.ctx, &emptypb.Empty{})
		if err != nil {
			panic(err)
		}
		c.typeReg = NewTypeRegistry()
		if err := c.typeReg.UseFileDescriptorSet(schema); err != nil {
			panic(err)
		}
	}
	return c.typeReg
}

func (c *Client) GetTypeNames() (map[string][]string, error) {
	rsp, err := c.cli.GetTypeNames(c.ctx, &emptypb.Empty{})
	res := map[string][]string{}
	if err != nil {
		return nil, err
	}
	for _, v := range rsp.TypeAliases {
		res[v.Fullname] = v.Aliases
	}
	return res, nil
}

func (c *Client) SearchRef(tk TypeKey, selector func(obj *ObjRef) bool) (chan *ObjRef, error) {
	ch := make(chan *ObjRef, 10)
	req := &StreamRefReq{
		TypeKey: tk[:],
	}
	stream, err := c.cli.StreamRefs(c.ctx, req)
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(ch)
		for {
			ref, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.StreamRefs(_) = _, %v", c.cli, err)
			}
			if selector(ref) {
				ch <- ref
			}
		}
	}()
	return ch, nil
}
