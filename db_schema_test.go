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
	"testing"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestDbSchema(t *testing.T) {
	tmpDir := CreateTestDb()
	defer CloseTestDb(tmpDir)

	r := NewTypeRegistry()
	if err := r.refresh(); err != nil {
		t.FailNow()
	}
	err := r.StoreSchema()
	if err != nil {
		t.FailNow()
	}
	r2 := NewTypeRegistry()
	err = r2.LoadSchema()
	if err != nil {
		t.FailNow()
	}

	fds1 := r.GetFileDescriptorSet()
	fds2 := r2.GetFileDescriptorSet()

	f1, err := protodesc.NewFiles(fds1)
	if err != nil {
		t.FailNow()
	}
	f2, err := protodesc.NewFiles(fds2)
	if err != nil {
		t.FailNow()
	}
	md1 := []protoreflect.MessageDescriptor{}
	ed1 := []protoreflect.EnumDescriptor{}
	f1.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Messages().Len(); i++ {
			md1 = append(md1, fd.Messages().Get(i))
		}
		for i := 0; i < fd.Enums().Len(); i++ {
			ed1 = append(ed1, fd.Enums().Get(i))
		}
		return true
	})
	md2 := []protoreflect.MessageDescriptor{}
	ed2 := []protoreflect.EnumDescriptor{}
	f2.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Messages().Len(); i++ {
			md2 = append(md2, fd.Messages().Get(i))
		}
		for i := 0; i < fd.Enums().Len(); i++ {
			ed2 = append(ed2, fd.Enums().Get(i))
		}
		return true
	})
loop2:
	for _, v1 := range md1 {
		for _, v2 := range md2 {
			if v1.FullName() == v2.FullName() {
				continue loop2
			}
		}
		t.FailNow() // Not same
	}
loop3:
	for _, v1 := range ed1 {
		for _, v2 := range ed2 {
			if v1.FullName() == v2.FullName() {
				continue loop3
			}
		}
		t.FailNow() // Not same
	}
}
