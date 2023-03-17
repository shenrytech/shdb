package shdb

import (
	"bytes"
	"net/url"

	"github.com/google/uuid"
)

type DbKey struct {
	data [20]byte
}

func (k *DbKey) Equal(other *DbKey) bool {
	return bytes.Equal(k.data[:], other.data[:])
}

func (k *DbKey) String() string {
	return url.QueryEscape(string(k.data[:]))
}

func DbKeyFromString(str string) (*DbKey, error) {
	k := &DbKey{}
	unescaped, err := url.QueryUnescape(str)
	if err != nil {
		return nil, err
	}
	k.data = [20]byte([]byte(unescaped)[:])
	return k, nil
}

func (k *DbKey) Key() []byte {
	return k.data[:]
}

func NewDbKey(typeKey TypeKey, id []byte) *DbKey {
	ret := &DbKey{}
	ret.SetType(typeKey)
	ret.SetUuidBytes(id)
	return ret
}

func MarshalDbKey(key []byte) *DbKey {
	ret := &DbKey{}
	copy(ret.data[:], key)
	return ret
}

func GetDbKey(obj IObject) *DbKey {
	ret := &DbKey{}
	ret.SetType([4]byte(obj.GetMetadata().Type))
	ret.SetUuidBytes(obj.GetMetadata().Uuid)
	return ret
}

func (k *DbKey) Uuid() uuid.UUID {
	id, err := uuid.FromBytes(k.data[4:])
	if err != nil {
		panic(err)
	}
	return id
}

func (k *DbKey) UuidBytes() []byte {
	return k.data[4:]
}

func (k *DbKey) TypeKey() TypeKey {
	return [4]byte(k.data[:4])
}

func (k *DbKey) SetUuidBytes(id []byte) {
	copy(k.data[4:], id[:])
}

func (k *DbKey) SetUuid(id uuid.UUID) {
	v, err := id.MarshalBinary()
	if err != nil {
		panic(err)
	}
	copy(k.data[4:], v[:])
}

func (k *DbKey) SetType(keyType TypeKey) {
	copy(k.data[:4], keyType[:])
}
