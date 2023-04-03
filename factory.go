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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// IObject is the interface for all objects in the dataabase.
// They all should stem from a protobuf message that looks like this:
//
//		message TObject {
//	  		shdb.Metadata metadata = 1;
//	  		string my_field = 2;
//	  		uint64 my_int = 3;
//		}
type IObject interface {
	proto.Message
	GetMetadata() *Metadata
}

// KeyVal is the binary representation of an IObject as it is stored in the database
type KeyVal struct {
	TypeId
	Value []byte
}

// TypeKey is the four bytes that identifies the type of an object
type TypeKey = [4]byte

var TypeKeyAll = TypeKey{0xff, 0xff, 0xff, 0xff}

// New creates a new IObject based on the type key and initializes
// the Metadata fields.
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

// MustNew is like `New` but panics if there is an error
func MustNew[T IObject](typeKey TypeKey) T {
	obj, err := New[T](typeKey)
	if err != nil {
		panic(err)
	}
	return obj
}

func create(typeKey TypeKey) (proto.Message, error) {
	obj, err := typeRegistry.CreateObject(typeKey)
	if err != nil {
		return nil, ErrNotAnObject
	}
	return obj, err
}

// Create just creates the memory for an IObject without
// initializing the Metadata
func Create[T IObject](typeKey TypeKey) (t T, err error) {

	obj, err := typeRegistry.CreateObject(typeKey)
	if err != nil {
		return t, ErrNotAnObject
	}
	return obj.(T), err
}

func unmarshal(kv KeyVal) (proto.Message, error) {
	obj, err := create(kv.TypeKey())
	if err != nil {
		return nil, err
	}
	if err = proto.Unmarshal(kv.Value, obj); err != nil {
		return nil, err
	}
	return obj, err
}

// Unmarshal returns the IObject from a KeyVal binary representation
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

// UnmarshalMany unmarshals a list of KeyVal binary representations
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

// Marshal returns a list of KeyVal binary representation representing
// a list of IObjects
func Marshal[T IObject](objs ...T) (ret []KeyVal, err error) {
	if len(objs) == 0 {
		return nil, nil
	}
	tk := objs[0].GetMetadata().TypeId().TypeKey()
	ret = []KeyVal{}
	for _, o := range objs {
		kv := KeyVal{}
		kv.SetType(tk)
		kv.SetUuidBytes(o.GetMetadata().Uuid)
		kv.Value, err = proto.Marshal(o)
		if err != nil {
			return nil, err
		}
		ret = append(ret, kv)
	}
	return ret, nil
}
