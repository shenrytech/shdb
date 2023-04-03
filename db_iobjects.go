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
	"context"
	"errors"
	"io"
	"log"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Put creates objects in the database. If the object already existed
// it will be overwritten.
func Put[T IObject](val ...T) error {
	if len(val) == 0 {
		return nil
	}
	for _, v := range val {
		v.GetMetadata().UpdatedAt = timestamppb.Now()
	}
	kv, err := Marshal(val...)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket_obj)
		for _, v := range kv {
			err = b.Put(v.Key(), v.Value)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, v := range val {
		notifyCreate(v)
	}
	return nil
}

func get(tid TypeId) (*KeyVal, error) {
	kv := &KeyVal{TypeId: tid}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket_obj)
		kv.Value = b.Get(kv.Key())
		if kv.Value == nil {
			return ErrNotFound
		}
		return nil
	})
	return kv, err
}

// Get an object from the database based on the type and id of the object.
func Get[T IObject](tid TypeId) (T, error) {
	var t T
	kv := KeyVal{TypeId: tid}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket_obj)
		kv.Value = b.Get(kv.Key())
		if kv.Value == nil {
			return ErrNotFound
		}
		var err error
		t, err = Unmarshal[T](kv)
		return err
	})
	return t, err
}

// GetRef returns an object from the database based on an ObjRef
func GetRef[T IObject](ref *ObjRef) (T, error) {
	var t T
	kv := KeyVal{TypeId: *ref.TypeId()}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket_obj)
		kv.Value = b.Get(kv.Key())
		if kv.Value == nil {
			return ErrNotFound
		}
		var err error
		t, err = Unmarshal[T](kv)
		return err
	})
	return t, err
}

// GetOne returns one of the objects in the database with the specified type that
// matches the selector function.
// The selector should return true for a match and false otherwise when presented
// with an object.
func GetFirst[T IObject](typeKey TypeKey, selector func(obj T) bool) (T, error) {
	var t T
	err := db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bucket_obj).Cursor()
		for k, v := c.Seek(typeKey[:]); k != nil && bytes.HasPrefix(k, typeKey[:]); k, v = c.Next() {
			kv := KeyVal{TypeId: *MarshalTypeId(k), Value: v}
			obj, err := Unmarshal[T](kv)
			if err != nil {
				return err
			}
			if selector(obj) {
				t = obj
				return nil
			}
		}
		return ErrNotFound
	})
	return t, err
}

// Update an object in the database by using an updater function. The updated
// object is returned.
func Update[T IObject](tid TypeId, updater func(obj T) (T, error)) (t T, err error) {
	var (
		prev T
		obj  T
	)

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket_obj)
		kv := KeyVal{TypeId: tid}
		kv.Value = b.Get(kv.Key())
		if kv.Value == nil {
			return ErrNotFound
		}
		var err error
		prev, err = Unmarshal[T](kv)
		if err != nil {
			return err
		}
		obj, err = updater(proto.Clone(prev).(T))
		if err != nil {
			return err
		}
		obj.GetMetadata().UpdatedAt = timestamppb.Now()
		kvs, err := Marshal(obj)
		b.Put(kvs[0].Key(), kvs[0].Value)
		return err
	})
	if err == nil {
		notifyUpdate(obj, prev)
	}
	return obj, err
}

// Delete an object from the database based on the type and id.
// The old value is returned
func Delete[T IObject](tid TypeId) (T, error) {
	obj, err := Get[T](tid)
	if err != nil {
		return obj, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket_obj)
		return b.Delete(tid.Key())
	})
	if err == nil {
		notifyDelete(obj)
	}
	return obj, err

}

// Delete all objects of a specific type from the database.
func DeleteAll(tk TypeKey) error {
	deleted := []IObject{}
	err := db.Update(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bucket_obj).Cursor()
		for k, v := c.Seek(tk[:]); k != nil && bytes.HasPrefix(k, tk[:]); k, v = c.Next() {
			kv := KeyVal{TypeId: *MarshalTypeId(k), Value: v}
			t, err := Unmarshal[IObject](kv)
			if err != nil {
				deleted = append(deleted, t)
			}
		}
		return nil
	})
	for _, d := range deleted {
		notifyDelete(d)
	}
	return err

}

// GetAllKV returns all KeyVals of the database.
func GetAllKV(typeKey TypeKey) ([]KeyVal, error) {
	allKvs := []KeyVal{}
	err := db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(bucket_obj).Cursor()
		for k, v := c.Seek(typeKey[:]); k != nil && bytes.HasPrefix(k, typeKey[:]); k, v = c.Next() {
			kv := KeyVal{TypeId: *MarshalTypeId(k), Value: v}
			allKvs = append(allKvs, kv)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// sort.SliceStable(allKvs, func(i, j int) bool {
	// 	return bytes.Compare(allKvs[i].Key(), allKvs[j].Key()) < 0
	// })
	return allKvs, nil
}

// GetAll returns all objects in database of a specific type
func GetAll[T IObject](typeKey TypeKey) ([]T, error) {
	allKvs, err := GetAllKV(typeKey)
	if err != nil {
		return nil, err
	}
	return UnmarshalMany[T](allKvs)
}

func queryStream(typ TypeKey, doneCh chan struct{}) (ch chan proto.Message) {
	ch = make(chan proto.Message, 10)
	go func() {
		defer close(ch)
		err := db.View(func(tx *bbolt.Tx) error {
			cnt := 1
			c := tx.Bucket(bucket_obj).Cursor()
			for k, v := c.Seek(typ[:]); k != nil && bytes.HasPrefix(k, typ[:]); k, v = c.Next() {
				kv := KeyVal{TypeId: *MarshalTypeId(k), Value: v}
				if kv.Value == nil {
					log.Printf("empty value in database kv=[%s]\n", kv.String())
				}
				t, err := unmarshal(kv)
				if err != nil {
					log.Printf("failed to parse value in database kv=[%s], err=[%v]\n", kv.String(), err)
				} else {
					select {
					case ch <- t:
						cnt++
					case <-doneCh:
						return io.EOF
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("queryStream failed, err=[%v]\n", err)
		}
	}()
	return
}

type activeStream struct {
	inCh   chan proto.Message
	doneCh chan struct{}
}

var activeStreams = map[uuid.UUID]*activeStream{}

// Query returns all objects of a specific type matching a selector function.
// Paging of the results is implemented using a pageSize and a token. If there are more
// results available after pageSize items have been returned a non-empty nextPageToken is returned
// that can be used to retrieve a new page of results.
// If nextPageToken is the empty string, no more results are available.
func Query[T IObject](ctx context.Context, typ TypeKey, selectFn func(obj T) (bool, error), pageSize int32, pageToken string) (result []T, nextPageToken string, err error) {

	var (
		streamId uuid.UUID
		stream   *activeStream
		ok       bool
	)

	// Find out the active stream, or create a new
	if pageToken == "" {
		streamId, err = uuid.NewUUID()
		if err != nil {
			return
		}
		doneCh := make(chan struct{})
		inCh := queryStream(typ, doneCh)
		stream = &activeStream{inCh: inCh, doneCh: doneCh}
		activeStreams[streamId] = stream
	} else {
		streamId, err = uuid.Parse(pageToken)
		if err != nil {
			return
		}
		stream, ok = activeStreams[streamId]
		if !ok {
			return nil, "", ErrSessionInvalid
		}
	}

	// Collect results from the stream
	res := []T{}
	for i := 0; i < int(pageSize); i++ {
		select {
		case obj, ok := <-stream.inCh:
			if !ok {
				close(stream.doneCh)
				delete(activeStreams, streamId)
				return res, "", nil
			}
			t := obj.(T)
			selected, err := selectFn(t)
			if selected {
				res = append(res, obj.(T))
			}
			if errors.Is(err, io.EOF) {
				close(stream.doneCh)
				delete(activeStreams, streamId)
				return res, "", nil
			}
		case <-ctx.Done():
			close(stream.doneCh)
			delete(activeStreams, streamId)
			return res, "", ErrContextCancelled
		}
	}

	return res, streamId.String(), nil
}

// List all objects pertaining to a specific type. For arguments and paging see `Query` method.
func List[T IObject](ctx context.Context, typ TypeKey, pageSize int32, pageToken string) (result []T, nextPageToken string, err error) {
	identityFn := func(a T) (bool, error) {
		return true, nil
	}
	return Query(ctx, typ, identityFn, pageSize, pageToken)
}
