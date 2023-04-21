// Copyright 2023 Shenry Tech AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Todo: Clean up Marshalling/Unmarshalling of ObjRef
// Todo: Remove either ObjRef or TypeID. They do the same thing

package shdb

import (
	"bytes"

	"github.com/google/uuid"
)

// Marshal an ObjRef to a byte slice
func (r *ObjRef) Marshal() []byte {
	res := []byte{}
	res = append(res, r.Type...)
	res = append(res, r.Uuid...)
	return res
}

// UnmarshalObjRef a byte slice into ObjRef
func UnmarshalObjRef(data []byte) (*ObjRef, error) {
	if len(data) != 20 {
		return nil, ErrInvalidType
	}
	res := &ObjRef{}
	res.Type = data[:4]
	res.Uuid = data[4:]
	return res, nil
}

// TypeId returns the ObjRef as a *TypeId
func (r *ObjRef) TypeId() *TypeId {
	res := &TypeId{}
	if len(r.Type) != 4 {
		panic("invalid ObjRef")
	}
	if len(r.Uuid) != 16 {
		panic("invalid ObjRef")
	}
	copy(res.data[:4], r.Type)
	copy(res.data[4:], r.Uuid)
	return res
}

func (r *ObjRef) UUID() uuid.UUID {
	res, err := uuid.FromBytes(r.Uuid)
	if err != nil {
		panic(err)
	}
	return res
}

func ObjRefFromUUID(typeKey TypeKey, id string) (*ObjRef, error) {
	u, e := uuid.Parse(id)
	if e != nil {
		return nil, e
	}
	ub, e := u.MarshalBinary()
	if e != nil {
		return nil, e
	}
	res := &ObjRef{
		Type: typeKey[:],
		Uuid: ub,
	}
	return res, nil
}

func ParseObjRef(typeKey TypeKey, id interface{}) (*ObjRef, error) {
	objRef := &ObjRef{
		Type: typeKey[:],
	}
	var err error
	switch val := id.(type) {
	case uuid.UUID:
		objRef.Uuid, err = val.MarshalBinary()
		if err != nil {
			return nil, err
		}
	case string:
		uuid, err := uuid.Parse(val)
		if err != nil {
			return nil, err
		}
		objRef.Uuid, err = uuid.MarshalBinary()
		if err != nil {
			return nil, err
		}
	case []byte:
		if len(val) != 16 {
			return nil, ErrInvalidType
		}
		objRef.Uuid = bytes.Clone(val)
	}
	return objRef, nil
}

func MustParseObjRef(typeKey TypeKey, id interface{}) *ObjRef {
	o, err := ParseObjRef(typeKey, id)
	if err != nil {
		panic(err)
	}
	return o
}

func (r *ObjRef) Equal(other *ObjRef) bool {
	if other == nil {
		return false
	}
	if !bytes.Equal(r.Type, other.Type) {
		return false
	}
	return bytes.Equal(r.Uuid, other.Uuid)
}
