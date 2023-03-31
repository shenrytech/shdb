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
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	schemaKey = []byte{1, 2, 3}
)

// StoreSchema stores the current state of a protoregistry.Files object in the
// schema bucket.
func StoreSchema(files protoregistry.Files) error {
	fileSet := &descriptorpb.FileDescriptorSet{}
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		fileSet.File = append(fileSet.File, protodesc.ToFileDescriptorProto(fd))
		return true
	})

	data, err := proto.Marshal(fileSet)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket_schema)
		b.Put(schemaKey, data)
		return nil
	})
}

func LoadSchema() (*descriptorpb.FileDescriptorSet, error) {
	data := []byte{}
	fds := &descriptorpb.FileDescriptorSet{}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket_schema)
		payload := b.Get(schemaKey)
		if payload == nil {
			return ErrNotFound
		}
		copy(data, payload)
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(data, fds)
	if err != nil {
		return nil, err
	}
	return fds, err
}
