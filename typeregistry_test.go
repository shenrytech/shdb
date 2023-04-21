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
	"log"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
)

func CreateTestFileDescriptor() *descriptorpb.FileDescriptorProto {
	sp := func(str string) *string { return &str }
	a := protodesc.ToFileDescriptorProto(File_pb_shdb_v1_test_proto)
	b := proto.Clone(a).(*descriptorpb.FileDescriptorProto)

	b.Package = sp("shdb.test.volatile")
	b.MessageType[0].Name = sp("TObj2")
	b.Name = sp("shdb/test/shdb_test_volatile.proto")
	return b
}

func TestTR(t *testing.T) {
	TObject := TypeKeyOf("shdb.v1.TObject")
	r := NewTypeRegistry()
	if err := r.refresh(); err != nil {
		t.FailNow()
	}

	obj, err := r.CreateObject(TObject)
	if err != nil {
		t.FailNow()
	}

	// Test that CreateObject generated a Metadata field
	if obj.GetMetadata() == nil {
		t.FailNow()
	}
	log.Printf("obj : %v", obj)

}

func TestTRDyn(t *testing.T) {
	TObj2 := TypeKeyOf("shdb.test.volatile.TObj2")
	r := NewTypeRegistry()
	if err := r.refresh(); err != nil {
		t.FailNow()
	}

	// Add the synthetic file and message
	fd := CreateTestFileDescriptor()
	if err := r.AddFileFromProtoFileDescriptor(fd); err != nil {
		t.FailNow()
	}
	// Create a new TObjs with metadata
	obj, err := r.CreateObject(TObj2)
	if err != nil {
		t.FailNow()
	}

	// Test that CreateObject generated a Metadata field
	if md := obj.GetMetadata(); md == nil {
		t.FailNow()
	} else {
		if md.CreatedAt == nil {
			t.Fail()
		}
		if md.Labels == nil {
			t.Fail()
		}
		if md.Type == nil {
			t.Fail()
		}
		if md.Uuid == nil {
			t.Fail()
		}
	}
}

func TestGetNames(t *testing.T) {
	r := NewTypeRegistry()
	if err := r.refresh(); err != nil {
		t.FailNow()
	}

	// Add the synthetic file and message
	fd := CreateTestFileDescriptor()
	if err := r.AddFileFromProtoFileDescriptor(fd); err != nil {
		t.FailNow()
	}
	nameAliases := r.GetTypeNames()
	if len(nameAliases) != 3 {
		t.FailNow()
	}
	if aliases1, ok := nameAliases["shdb.v1.TObject"]; !ok {
		t.FailNow()
	} else {
		if len(aliases1) != 2 {
			t.FailNow()
		}
	}
	if aliases2, ok := nameAliases["shdb.test.volatile.TObj2"]; !ok {
		t.FailNow()
	} else {
		if len(aliases2) != 2 {
			t.FailNow()
		}

	}
}
