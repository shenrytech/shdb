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
	"github.com/shenrytech/shdb/jsonsearch"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// SearchProto searches within the fields of the proto message forthe string provided.
// The result will contain a number of hits in the form
// - /[fieldName|@idx]/...
// Examples:
//   - / - The object contained only one item and it matched
//   - /myField - {"myField": <match>}
//   - /field1/field2/@3 {"field1": {"field2": [1,2,<match>]}}
func SearchProto(m proto.Message, query func(s string) bool) (hits []string, err error) {
	mo := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	jsonData, err := mo.Marshal(m)
	if err != nil {
		return nil, err
	}
	p := jsonsearch.NewParser(jsonData, query)
	err = p.Parse("")
	return p.FieldPaths, err
}

func searchStream(typ TypeKey, selector func(s string) bool, doneCh chan struct{}) (ch chan *SearchHit) {
	ch = make(chan *SearchHit, 10)
	go func() {
		defer func() {
			close(ch)
		}()

		err := db.View(func(tx *bbolt.Tx) error {
			cnt := 1
			c := tx.Bucket(bucket_obj).Cursor()
			for k, v := c.Seek(typ[:]); k != nil && bytes.HasPrefix(k, typ[:]); k, v = c.Next() {
				kv := KeyVal{TypeId: *MarshalTypeId(k), Value: v}
				if kv.Value == nil {
					log.Printf("empty value in database kv=[%s]\n", kv.String())
					continue
				}
				t, err := Unmarshal[IObject](kv)
				if err != nil {
					log.Printf("failed to parse value in database kv=[%s], err=[%v]\n", kv.String(), err)
				} else {
					hits, err := SearchProto(t, selector)
					if err == nil && len(hits) > 0 {
						select {
						case ch <- &SearchHit{
							Hits:     hits,
							Metadata: t.GetMetadata(),
						}:
							cnt++
						case <-doneCh:
							return io.EOF
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("searchStream failed, err=[%v]\n", err)
		}
	}()
	return
}

type activeSearchStream struct {
	inCh   chan *SearchHit
	doneCh chan struct{}
}

var activeSearchStreams = map[uuid.UUID]*activeSearchStream{}

// Search searches the values of the fields of objects pertaining to a type by calling
// a selector function for each field in all objects.
// For paging functionality see `Query` method.
func Search(ctx context.Context,
	typ TypeKey,
	selector func(string) bool,
	pageSize int32,
	pageToken string) (result *SearchResult, nextPageToken string, err error) {

	var (
		streamId uuid.UUID
		stream   *activeSearchStream
		ok       bool
	)

	// Find out the active stream, or create a new
	if pageToken == "" {
		streamId, err = uuid.NewUUID()
		if err != nil {
			return
		}
		doneCh := make(chan struct{})
		inCh := searchStream(typ, selector, doneCh)
		stream = &activeSearchStream{inCh: inCh, doneCh: doneCh}
		activeSearchStreams[streamId] = stream
	} else {
		streamId, err = uuid.Parse(pageToken)
		if err != nil {
			return
		}
		stream, ok = activeSearchStreams[streamId]
		if !ok {
			return nil, "", ErrSessionInvalid
		}
	}

	// Collect results from the stream
	res := &SearchResult{
		Hits: []*SearchHit{},
	}
	for i := 0; i < int(pageSize); i++ {
		select {
		case obj, ok := <-stream.inCh:
			if !ok {
				close(stream.doneCh)
				delete(activeStreams, streamId)
				return res, "", nil
			}
			res.Hits = append(res.Hits, obj)
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
