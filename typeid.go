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

package shdb

import (
	"bytes"
	"net/url"

	"github.com/google/uuid"
)

// TypeId is the key in the key-value database that is used
// to store IObjects. It has a memory layout like:
//
//	[b0 .. b3]		TypeKey
//	[b4 .. b20]		Binary representation of an UUID
type TypeId struct {
	data [20]byte
}

// Equal compares two TypeIds and return true if they are equal
func (k *TypeId) Equal(other *TypeId) bool {
	return bytes.Equal(k.data[:], other.data[:])
}

// String returns the URL-encoded string of the TypeId
func (k *TypeId) String() string {
	return url.QueryEscape(string(k.data[:]))
}

// TypeIdFromString returns a TypeId from an URL-encoded string
func TypeIdFromString(str string) (*TypeId, error) {
	k := &TypeId{}
	unescaped, err := url.QueryUnescape(str)
	if err != nil {
		return nil, err
	}
	k.data = [20]byte([]byte(unescaped)[:])
	return k, nil
}

// Key returns the TypeId as a []byte slice
func (k *TypeId) Key() []byte {
	return k.data[:]
}

// NewTypeId creates a new TypeId based on TypeKey and UUID (Byte version)
func NewTypeId(typeKey TypeKey, id []byte) *TypeId {
	ret := &TypeId{}
	ret.SetType(typeKey)
	ret.SetUuidBytes(id)
	return ret
}

// MarshalTypeId creates a TypeId from a []byte slice
func MarshalTypeId(data []byte) *TypeId {
	ret := &TypeId{}
	copy(ret.data[:], data)
	return ret
}

// GetTypeId returns the TypeID from the Metadata of an IObject
func GetTypeId(obj IObject) *TypeId {
	ret := &TypeId{}
	ret.SetType([4]byte(obj.GetMetadata().Type))
	ret.SetUuidBytes(obj.GetMetadata().Uuid)
	return ret
}

// Uuid returns the uuid.UUID of the id of a TypeID
func (k TypeId) Uuid() uuid.UUID {
	id, err := uuid.FromBytes(k.data[4:])
	if err != nil {
		panic(err)
	}
	return id
}

// UuidBytes returns the byte version of the id of a TypeId
func (k TypeId) UuidBytes() []byte {
	return k.data[4:]
}

// TypeKey returns the TypeKey of a TypeId
func (k TypeId) TypeKey() TypeKey {
	return [4]byte(k.data[:4])
}

// SetUuidBytes sets the Id (as bytes) of a TypeId
func (k *TypeId) SetUuidBytes(id []byte) {
	copy(k.data[4:], id[:])
}

// SetUuid sets thte id (as uuid.UUID) of a TypeId
func (k *TypeId) SetUuid(id uuid.UUID) {
	v, err := id.MarshalBinary()
	if err != nil {
		panic(err)
	}
	copy(k.data[4:], v[:])
}

// SetType sets the type of a TypeId
func (k *TypeId) SetType(keyType TypeKey) {
	copy(k.data[:4], keyType[:])
}
