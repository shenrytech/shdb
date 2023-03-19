//
//Copyright 2023 (C) Shenry Tech AB
//
//License: MIT
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.2
// source: shdb.proto

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

type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        []byte                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Uuid        []byte                 `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Labels      []string               `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty"`
	Description string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shdb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_shdb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_shdb_proto_rawDescGZIP(), []int{0}
}

func (x *Metadata) GetType() []byte {
	if x != nil {
		return x.Type
	}
	return nil
}

func (x *Metadata) GetUuid() []byte {
	if x != nil {
		return x.Uuid
	}
	return nil
}

func (x *Metadata) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Metadata) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Metadata) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Metadata) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type SearchHit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Hits     []string  `protobuf:"bytes,2,rep,name=hits,proto3" json:"hits,omitempty"`
}

func (x *SearchHit) Reset() {
	*x = SearchHit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shdb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchHit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchHit) ProtoMessage() {}

func (x *SearchHit) ProtoReflect() protoreflect.Message {
	mi := &file_shdb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchHit.ProtoReflect.Descriptor instead.
func (*SearchHit) Descriptor() ([]byte, []int) {
	return file_shdb_proto_rawDescGZIP(), []int{1}
}

func (x *SearchHit) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *SearchHit) GetHits() []string {
	if x != nil {
		return x.Hits
	}
	return nil
}

type SearchResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hits []*SearchHit `protobuf:"bytes,1,rep,name=hits,proto3" json:"hits,omitempty"`
}

func (x *SearchResult) Reset() {
	*x = SearchResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shdb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchResult) ProtoMessage() {}

func (x *SearchResult) ProtoReflect() protoreflect.Message {
	mi := &file_shdb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchResult.ProtoReflect.Descriptor instead.
func (*SearchResult) Descriptor() ([]byte, []int) {
	return file_shdb_proto_rawDescGZIP(), []int{2}
}

func (x *SearchResult) GetHits() []*SearchHit {
	if x != nil {
		return x.Hits
	}
	return nil
}

type TObject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	MyField  string    `protobuf:"bytes,2,opt,name=my_field,json=myField,proto3" json:"my_field,omitempty"`
	MyInt    uint64    `protobuf:"varint,3,opt,name=my_int,json=myInt,proto3" json:"my_int,omitempty"`
}

func (x *TObject) Reset() {
	*x = TObject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shdb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TObject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TObject) ProtoMessage() {}

func (x *TObject) ProtoReflect() protoreflect.Message {
	mi := &file_shdb_proto_msgTypes[3]
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
	return file_shdb_proto_rawDescGZIP(), []int{3}
}

func (x *TObject) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *TObject) GetMyField() string {
	if x != nil {
		return x.MyField
	}
	return ""
}

func (x *TObject) GetMyInt() uint64 {
	if x != nil {
		return x.MyInt
	}
	return 0
}

var File_shdb_proto protoreflect.FileDescriptor

var file_shdb_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x73, 0x68, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x73, 0x68,
	0x64, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xe2, 0x01, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x4b, 0x0a, 0x09, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x48, 0x69, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x68, 0x64, 0x62, 0x2e, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x69, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x04, 0x68, 0x69, 0x74, 0x73, 0x22, 0x33, 0x0a, 0x0c, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x23, 0x0a, 0x04, 0x68, 0x69, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x68, 0x64, 0x62, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x48, 0x69, 0x74, 0x52, 0x04, 0x68, 0x69, 0x74, 0x73, 0x22, 0x67, 0x0a, 0x07, 0x54, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x68, 0x64, 0x62, 0x2e, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x79, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x79, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x15, 0x0a, 0x06,
	0x6d, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6d, 0x79,
	0x49, 0x6e, 0x74, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x68, 0x65, 0x6e, 0x72, 0x79, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x73, 0x68, 0x64,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shdb_proto_rawDescOnce sync.Once
	file_shdb_proto_rawDescData = file_shdb_proto_rawDesc
)

func file_shdb_proto_rawDescGZIP() []byte {
	file_shdb_proto_rawDescOnce.Do(func() {
		file_shdb_proto_rawDescData = protoimpl.X.CompressGZIP(file_shdb_proto_rawDescData)
	})
	return file_shdb_proto_rawDescData
}

var file_shdb_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_shdb_proto_goTypes = []interface{}{
	(*Metadata)(nil),              // 0: shdb.Metadata
	(*SearchHit)(nil),             // 1: shdb.SearchHit
	(*SearchResult)(nil),          // 2: shdb.SearchResult
	(*TObject)(nil),               // 3: shdb.TObject
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_shdb_proto_depIdxs = []int32{
	4, // 0: shdb.Metadata.created_at:type_name -> google.protobuf.Timestamp
	4, // 1: shdb.Metadata.updated_at:type_name -> google.protobuf.Timestamp
	0, // 2: shdb.SearchHit.metadata:type_name -> shdb.Metadata
	1, // 3: shdb.SearchResult.hits:type_name -> shdb.SearchHit
	0, // 4: shdb.TObject.metadata:type_name -> shdb.Metadata
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_shdb_proto_init() }
func file_shdb_proto_init() {
	if File_shdb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shdb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metadata); i {
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
		file_shdb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchHit); i {
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
		file_shdb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchResult); i {
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
		file_shdb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_shdb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shdb_proto_goTypes,
		DependencyIndexes: file_shdb_proto_depIdxs,
		MessageInfos:      file_shdb_proto_msgTypes,
	}.Build()
	File_shdb_proto = out.File
	file_shdb_proto_rawDesc = nil
	file_shdb_proto_goTypes = nil
	file_shdb_proto_depIdxs = nil
}
