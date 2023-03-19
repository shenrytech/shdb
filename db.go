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

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	BUCKET_OBJ = []byte("obj")
	db         *bbolt.DB
	log        *zap.Logger
)

func Init(logger *zap.Logger, dbFile string) {
	var err error
	db, err = bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists(BUCKET_OBJ)
		return nil
	})
	log = logger
}

func Close() {
	if err := db.Close(); err != nil {
		log.Error("error closing database", zap.Error(err))
	} else {
		log.Debug("closed database", zap.String("dbFile", db.Path()))
	}
}

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
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BUCKET_OBJ)
		for _, v := range kv {
			err = b.Put(v.Key(), v.Value)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func Get[T IObject](tid TypeId) (T, error) {
	var t T
	kv := KeyVal{TypeId: tid}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BUCKET_OBJ)
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

func Update[T IObject](tid TypeId, fn func(obj T) (T, error)) (t T, err error) {
	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BUCKET_OBJ)
		kv := KeyVal{TypeId: tid}
		kv.Value = b.Get(kv.Key())
		if kv.Value == nil {
			return ErrNotFound
		}
		wo, err := Unmarshal[T](kv)
		if err != nil {
			return err
		}
		t, err = fn(wo)
		t.GetMetadata().UpdatedAt = timestamppb.Now()
		return err
	})
	return t, err
}

func Delete[T IObject](tid TypeId) (T, error) {
	val, err := Get[T](tid)
	if err != nil {
		return val, err
	}
	return val, db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BUCKET_OBJ)
		return b.Delete(tid.Key())
	})

}
func GetAllKV(typeKey TypeKey) ([]KeyVal, error) {
	allKvs := []KeyVal{}
	err := db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(BUCKET_OBJ).Cursor()
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
			c := tx.Bucket(BUCKET_OBJ).Cursor()
			for k, v := c.Seek(typ[:]); k != nil && bytes.HasPrefix(k, typ[:]); k, v = c.Next() {
				kv := KeyVal{TypeId: *MarshalTypeId(k), Value: v}
				if kv.Value == nil {
					log.Warn("empty value in database", zap.String("kv", kv.String()))
				}
				t, err := unmarshal(kv)
				if err != nil {
					log.Error("failed to parse value in database", zap.String("kv", kv.String()), zap.Error(err))
				} else {
					select {
					case ch <- t:
						log.Debug("sending msg", zap.Int("count", cnt))
						cnt++
					case <-doneCh:
						log.Debug("doneCh")
						return io.EOF
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Error("QueryStream failed", zap.Error(err))
		}
	}()
	return
}

type activeStream struct {
	inCh   chan proto.Message
	doneCh chan struct{}
}

var activeStreams = map[uuid.UUID]*activeStream{}

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

func List[T IObject](ctx context.Context, typ TypeKey, pageSize int32, pageToken string) (result []T, nextPageToken string, err error) {
	identityFn := func(a T) (bool, error) {
		return true, nil
	}
	return Query(ctx, typ, identityFn, pageSize, pageToken)
}
