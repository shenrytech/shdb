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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.2
// source: test.proto

//
//
//Key:
//Type 32bit
//UUID 16*8bit
//Labels
//Description
//Timestamps
//MessagePayload

package shdb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TObject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata  *Metadata              `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	MyInt     uint64                 `protobuf:"varint,2,opt,name=my_int,json=myInt,proto3" json:"my_int,omitempty"`
	MyString  string                 `protobuf:"bytes,3,opt,name=my_string,json=myString,proto3" json:"my_string,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Any       *anypb.Any             `protobuf:"bytes,5,opt,name=any,proto3" json:"any,omitempty"`
}

func (x *TObject) Reset() {
	*x = TObject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TObject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TObject) ProtoMessage() {}

func (x *TObject) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TObject.ProtoReflect.Descriptor instead.
func (*TObject) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *TObject) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *TObject) GetMyInt() uint64 {
	if x != nil {
		return x.MyInt
	}
	return 0
}

func (x *TObject) GetMyString() string {
	if x != nil {
		return x.MyString
	}
	return ""
}

func (x *TObject) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *TObject) GetAny() *anypb.Any {
	if x != nil {
		return x.Any
	}
	return nil
}

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x73, 0x68,
	0x64, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0a,
	0x73, 0x68, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd0, 0x02, 0x0a, 0x07, 0x54,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x68, 0x64, 0x62, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x6d, 0x79, 0x49, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x79, 0x5f,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x79,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x26, 0x0a, 0x03, 0x61, 0x6e, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x41, 0x6e, 0x79, 0x52, 0x03, 0x61, 0x6e, 0x79, 0x3a, 0x82, 0x01, 0x82, 0xb2, 0x19, 0x7e, 0x0a,
	0x07, 0x54, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x04, 0x74, 0x6f, 0x62, 0x6a, 0x12, 0x05,
	0x6e, 0x69, 0x73, 0x73, 0x65, 0x1a, 0x02, 0x23, 0x24, 0x22, 0x18, 0x0a, 0x0b, 0x69, 0x74, 0x65,
	0x6d, 0x5f, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x09, 0x61, 0x73, 0x64, 0x61, 0x73, 0x64,
	0x61, 0x73, 0x64, 0x22, 0x16, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x66, 0x75, 0x6c, 0x6c,
	0x12, 0x09, 0x61, 0x73, 0x64, 0x61, 0x73, 0x64, 0x61, 0x73, 0x64, 0x22, 0x18, 0x0a, 0x0b, 0x6c,
	0x69, 0x73, 0x74, 0x5f, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x09, 0x61, 0x73, 0x64, 0x61,
	0x73, 0x64, 0x61, 0x73, 0x64, 0x22, 0x16, 0x0a, 0x09, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x66, 0x75,
	0x6c, 0x6c, 0x12, 0x09, 0x61, 0x73, 0x64, 0x61, 0x73, 0x64, 0x61, 0x73, 0x64, 0x42, 0x1c, 0x5a,
	0x1a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x68, 0x65, 0x6e,
	0x72, 0x79, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x73, 0x68, 0x64, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_test_proto_rawDescOnce sync.Once
	file_test_proto_rawDescData = file_test_proto_rawDesc
)

func file_test_proto_rawDescGZIP() []byte {
	file_test_proto_rawDescOnce.Do(func() {
		file_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_proto_rawDescData)
	})
	return file_test_proto_rawDescData
}

var file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_test_proto_goTypes = []interface{}{
	(*TObject)(nil),               // 0: shdb.TObject
	(*Metadata)(nil),              // 1: shdb.Metadata
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
	(*anypb.Any)(nil),             // 3: google.protobuf.Any
}
var file_test_proto_depIdxs = []int32{
	1, // 0: shdb.TObject.metadata:type_name -> shdb.Metadata
	2, // 1: shdb.TObject.timestamp:type_name -> google.protobuf.Timestamp
	3, // 2: shdb.TObject.any:type_name -> google.protobuf.Any
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_test_proto_init() }
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	file_shdb_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TObject); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
		MessageInfos:      file_test_proto_msgTypes,
	}.Build()
	File_test_proto = out.File
	file_test_proto_rawDesc = nil
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}
