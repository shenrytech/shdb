package shdb

import (
	"bytes"
	"net/url"

	"github.com/google/uuid"
)

type TypeId struct {
	data [20]byte
}

func (k *TypeId) Equal(other *TypeId) bool {
	return bytes.Equal(k.data[:], other.data[:])
}

func (k *TypeId) String() string {
	return url.QueryEscape(string(k.data[:]))
}

func TypeIdFromString(str string) (*TypeId, error) {
	k := &TypeId{}
	unescaped, err := url.QueryUnescape(str)
	if err != nil {
		return nil, err
	}
	k.data = [20]byte([]byte(unescaped)[:])
	return k, nil
}

func (k *TypeId) Key() []byte {
	return k.data[:]
}

func NewTypeId(typeKey TypeKey, id []byte) *TypeId {
	ret := &TypeId{}
	ret.SetType(typeKey)
	ret.SetUuidBytes(id)
	return ret
}

func MarshalTypeId(key []byte) *TypeId {
	ret := &TypeId{}
	copy(ret.data[:], key)
	return ret
}

func GetTypeId(obj IObject) *TypeId {
	ret := &TypeId{}
	ret.SetType([4]byte(obj.GetMetadata().Type))
	ret.SetUuidBytes(obj.GetMetadata().Uuid)
	return ret
}

func (k *TypeId) Uuid() uuid.UUID {
	id, err := uuid.FromBytes(k.data[4:])
	if err != nil {
		panic(err)
	}
	return id
}

func (k *TypeId) UuidBytes() []byte {
	return k.data[4:]
}

func (k *TypeId) TypeKey() TypeKey {
	return [4]byte(k.data[:4])
}

func (k *TypeId) SetUuidBytes(id []byte) {
	copy(k.data[4:], id[:])
}

func (k *TypeId) SetUuid(id uuid.UUID) {
	v, err := id.MarshalBinary()
	if err != nil {
		panic(err)
	}
	copy(k.data[4:], v[:])
}

func (k *TypeId) SetType(keyType TypeKey) {
	copy(k.data[:4], keyType[:])
}
