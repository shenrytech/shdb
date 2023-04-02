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
	"context"
	"errors"
	"io"
	"log"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

func searchRefStream(selector func(ref *ObjRef) bool, doneCh chan struct{}) (ch chan *ObjRef) {
	ch = make(chan *ObjRef, 10)
	go func() {
		defer func() {
			close(ch)
		}()

		err := db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket(bucket_obj)
			err := b.ForEach(func(k, v []byte) error {
				ref, err := UnmarshalObjRef(k)
				if err != nil {
					return err
				}
				if selector(ref) {
					select {
					case ch <- ref:
					case <-doneCh:
						return io.EOF
					}
				}
				return nil
			})
			return err
		})

		if err != nil {
			log.Printf("searchStream failed, err=[%v]\n", err)
		}
	}()
	return
}

type activeSearchRefStream struct {
	inCh   chan *ObjRef
	doneCh chan struct{}
}

var activeSearchRefStreams = map[uuid.UUID]*activeSearchRefStream{}

// SearchRef searches the Ref of objects
// For paging functionality see `Query` method.
func SearchRef(ctx context.Context,
	selector func(*ObjRef) bool,
	pageSize int32,
	pageToken string) (result []*ObjRef, nextPageToken string, err error) {

	var (
		streamId uuid.UUID
		stream   *activeSearchRefStream
		ok       bool
	)

	// Find out the active stream, or create a new
	if pageToken == "" {
		streamId, err = uuid.NewUUID()
		if err != nil {
			return
		}
		doneCh := make(chan struct{})
		inCh := searchRefStream(selector, doneCh)
		stream = &activeSearchRefStream{inCh: inCh, doneCh: doneCh}
		activeSearchRefStreams[streamId] = stream
	} else {
		streamId, err = uuid.Parse(pageToken)
		if err != nil {
			return
		}
		stream, ok = activeSearchRefStreams[streamId]
		if !ok {
			return nil, "", ErrSessionInvalid
		}
	}

	// Collect results from the stream
	res := []*ObjRef{}
	for i := 0; i < int(pageSize); i++ {
		select {
		case obj, ok := <-stream.inCh:
			if !ok {
				close(stream.doneCh)
				delete(activeStreams, streamId)
				return res, "", nil
			}
			res = append(res, obj)
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
