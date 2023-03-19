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
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IObject interface {
	proto.Message
	GetMetadata() *Metadata
}

type KeyVal struct {
	TypeId
	Value []byte
}

type TypeKey = [4]byte

type typeInfo struct {
	typeKey  TypeKey
	tmplVal  proto.Message
	fullName protoreflect.FullName
}

var typeRegistry = map[TypeKey]typeInfo{}
var fnIndex = map[protoreflect.FullName]TypeKey{}

func Register[T IObject](tmpl T) {
	typeKey := TypeKey(tmpl.GetMetadata().Type)
	typeRegistry[typeKey] = typeInfo{
		typeKey:  typeKey,
		tmplVal:  proto.Clone(tmpl),
		fullName: tmpl.ProtoReflect().Descriptor().FullName(),
	}
	fnIndex[tmpl.ProtoReflect().Descriptor().FullName()] = typeKey
}

func New[T IObject](typeKey TypeKey) (obj T, err error) {
	var id []byte
	obj, err = Create[T](typeKey)
	if err != nil {
		return
	}
	id, err = uuid.New().MarshalBinary()
	if err != nil {
		return
	}
	obj.GetMetadata().Type = typeKey[:]
	obj.GetMetadata().CreatedAt = timestamppb.Now()
	obj.GetMetadata().UpdatedAt = timestamppb.Now()
	obj.GetMetadata().Uuid = id
	obj.GetMetadata().Labels = make([]string, 0)

	return obj, nil
}

func MustNew[T IObject](typeKey TypeKey) T {
	obj, err := New[T](typeKey)
	if err != nil {
		panic(err)
	}
	return obj
}

func create(typeKey TypeKey) (proto.Message, error) {
	tInfo, ok := typeRegistry[typeKey]
	if !ok {
		return nil, ErrNotAnObject
	}
	obj := proto.Clone(tInfo.tmplVal)
	return obj, nil
}

func Create[T IObject](typeKey TypeKey) (T, error) {
	tInfo, ok := typeRegistry[typeKey]
	if !ok {
		var t T
		return t, ErrNotAnObject
	}
	obj := proto.Clone(tInfo.tmplVal)
	return obj.(T), nil
}

func unmarshal(kv KeyVal) (proto.Message, error) {
	obj, err := create(kv.TypeKey())
	if err != nil {
		return nil, err
	}
	obj.ProtoReflect().Descriptor().FullName()
	if err = proto.Unmarshal(kv.Value, obj); err != nil {
		return nil, err
	}
	return obj, err
}

func Unmarshal[T IObject](kv KeyVal) (T, error) {
	obj, err := Create[T](kv.TypeKey())
	if err != nil {
		var t T
		return t, err
	}
	obj.ProtoReflect().Descriptor().FullName()
	if err = proto.Unmarshal(kv.Value, obj); err != nil {
		var t T
		return t, err
	}
	return obj, err
}

func UnmarshalMany[T IObject](kvs []KeyVal) ([]T, error) {
	res := []T{}
	for _, v := range kvs {
		obj, err := Unmarshal[T](v)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}
	return res, nil
}

func Marshal[T IObject](objs ...T) (ret []KeyVal, err error) {
	if len(objs) == 0 {
		return nil, nil
	}
	fn := objs[0].ProtoReflect().Descriptor().FullName()
	typeKey, ok := fnIndex[fn]
	if !ok {
		return nil, ErrInvalidType
	}
	ret = []KeyVal{}

	for _, o := range objs {
		kv := KeyVal{}
		kv.SetType(typeKey)
		kv.SetUuidBytes(o.GetMetadata().Uuid)
		kv.Value, err = proto.Marshal(o)
		if err != nil {
			return nil, err
		}
		ret = append(ret, kv)
	}
	return ret, nil
}
