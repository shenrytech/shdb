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
	cli     BinaryObjectServiceClient
	ctx     context.Context
	typeReg *TypeRegistry
}

// NewClient returns a new client for use with the API
func NewClient(ctx context.Context, cc *grpc.ClientConn) *Client {
	return &Client{ctx: ctx, cc: cc, cli: NewBinaryObjectServiceClient(cc), typeReg: nil}
}

func (c *Client) Get(tid TypeId) (IObject, error) {
	ref, err := UnmarshalObjRef(tid.Key())
	if err != nil {
		return nil, err
	}
	o, err := c.cli.Get(c.ctx, &GetReq{Ref: ref})
	if err != nil {
		return nil, err
	}
	return c.TypeRegistry().Unmarshal(o.Key, o.Value)
}

func (c *Client) List(tk TypeKey) ([]IObject, error) {
	// Let's not complicate things. Get only first 100000 objects...
	req := &ListReq{
		Type:      tk[:],
		PageSize:  100000,
		PageToken: "",
	}
	rsp, err := c.cli.List(c.ctx, req)
	if err != nil {
		return nil, err
	}
	res := []IObject{}
	for _, v := range rsp.Items {
		obj, err := c.TypeRegistry().Unmarshal(v.Key, v.Value)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}
	return res, nil
}

func (c *Client) Delete(tid TypeId) (IObject, error) {
	ref, err := UnmarshalObjRef(tid.Key())
	if err != nil {
		return nil, err
	}
	rsp, err := c.cli.Delete(c.ctx, &DeleteReq{Ref: ref})
	if err != nil {
		return nil, err
	}
	return c.TypeRegistry().Unmarshal(rsp.Key, rsp.Value)
}

func (c *Client) Create(typ TypeKey) (IObject, error) {
	rsp, err := c.cli.Create(c.ctx, &CreateReq{Type: typ[:]})
	if err != nil {
		return nil, err
	}
	return c.TypeRegistry().Unmarshal(rsp.Key, rsp.Value)
}

func (c *Client) Update(obj IObject) (IObject, error) {
	kvs, err := Marshal(obj)
	if err != nil {
		return nil, err
	}
	o := &BinaryObject{Key: kvs[0].Key(), Value: kvs[0].Value}
	rsp, err := c.cli.Update(c.ctx, &UpdateReq{Item: o})
	if err != nil {
		return nil, err
	}
	return c.TypeRegistry().Unmarshal(rsp.Key, rsp.Value)
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
