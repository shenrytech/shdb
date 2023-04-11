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
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TypeKeyOf(fullname string) TypeKey {
	h := fnv.New32a()
	h.Write([]byte(fullname))
	return TypeKey(h.Sum(nil))
}

type MessageInfo struct {
	Fullname       string
	Aliases        []string
	PrintTemplates map[string]string
	TypeKey        TypeKey
	MessageType    protoreflect.MessageType
	IsDynamic      bool
}

type TypeRegistry struct {
	fromFullname map[string]*MessageInfo
	fromTypeKey  map[TypeKey]*MessageInfo

	types *protoregistry.Types
	files *protoregistry.Files

	mux *sync.Mutex
}

func NewTypeRegistry() *TypeRegistry {

	r := &TypeRegistry{
		mux:          new(sync.Mutex),
		fromFullname: map[string]*MessageInfo{},
		fromTypeKey:  map[TypeKey]*MessageInfo{},
		types:        protoregistry.GlobalTypes,
		files:        protoregistry.GlobalFiles,
	}
	r.refresh()
	log.Println("Types registered:")
	for _, v := range r.fromFullname {
		log.Printf("%s: [%v]\n", v.Fullname, v.Aliases)
	}
	return r
}

func (r *TypeRegistry) clear() {
	r.fromFullname = nil
	r.fromTypeKey = nil
	r.fromFullname = map[string]*MessageInfo{}
	r.fromTypeKey = map[TypeKey]*MessageInfo{}
}

func (r *TypeRegistry) addMessage(md protoreflect.MessageDescriptor) error {
	mi := &MessageInfo{
		Fullname:       string(md.FullName()),
		TypeKey:        TypeKeyOf(string(md.FullName())),
		PrintTemplates: map[string]string{},
		Aliases:        []string{},
		IsDynamic:      false,
	}

	// Create a protoreflect.MessageType that can be used to
	// create instances of the object. If there is a Go type use
	// that, otherwise it's a dynamic type
	mt, err := r.types.FindMessageByName(md.FullName())
	if err != nil {
		if errors.Is(err, protoregistry.NotFound) {
			mt = dynamicpb.NewMessageType(md)
			mi.IsDynamic = true

		} else {
			return err
		}
	}
	mi.MessageType = mt

	if proto.HasExtension(md.Options(), E_ShdbOptions) {
		ext := proto.GetExtension(md.Options(), E_ShdbOptions).(*Shdb_Message_Options)
		if ext == nil {
			panic("Shdb_Options is of invalid type")
		}
		for k, v := range ext.PrintTemplates {
			mi.PrintTemplates[k] = v
		}
		mi.Aliases = append(mi.Aliases, ext.Aliases...)
		// mi.TypeKey = TypeKey(ext.TypeKey) - TypeKey is now from hashing the fullname
	}
	r.fromFullname[mi.Fullname] = mi
	r.fromTypeKey[[4]byte(mi.TypeKey)] = mi
	return nil
}

func (r *TypeRegistry) addFile(fd protoreflect.FileDescriptor) error {
	_, err := r.files.FindFileByPath(fd.Path())
	if err != nil {
		if !errors.Is(err, protoregistry.NotFound) {
			return err
		} else {
			if err := r.files.RegisterFile(fd); err != nil {
				log.Printf("error registerung proto file %s already registered [%v]", fd.Path(), err)
			} else {
				log.Printf("adding %s to fds", fd.Path())
			}
		}
	}
	return nil
}

func (r *TypeRegistry) AddFile(fd protoreflect.FileDescriptor) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if err := r.addFile(fd); err != nil {
		return err
	}
	return r.refresh()
}

func (r *TypeRegistry) AddFileFromProtoFileDescriptor(fd *descriptorpb.FileDescriptorProto) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	o := protodesc.FileOptions{}
	pfd, err := o.New(fd, r.files)
	if err != nil {
		return err
	}
	if err := r.addFile(pfd); err != nil {
		return err
	}
	return r.refresh()
}

func (r *TypeRegistry) refresh() error {

	var err error

	r.files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for idx := 0; idx < fd.Messages().Len(); idx++ {
			md := fd.Messages().Get(idx)
			// Make sure md has a field named 'metadata' and that it is of type
			// shdb.Metadata
			fd := md.Fields().ByName("metadata")
			if fd == nil {
				return true
			}
			fmd := fd.Message()
			if fmd == nil {
				return true
			}
			if fmd.FullName() != "shdb.v1.Metadata" {
				return true
			}
			// Include message if it's not already there
			if _, ok := r.fromFullname[string(md.FullName())]; !ok {
				if err = r.addMessage(md); err != nil {
					return false
				}
			}
		}
		return true
	})
	return err
}

func (r *TypeRegistry) CreateEmptyObject(tk TypeKey) (IObject, error) {
	mi, ok := r.fromTypeKey[tk]
	if !ok {
		return nil, ErrNotFound

	}
	var obj IObject
	if mi.IsDynamic {
		obj = &DynObject{Message: mi.MessageType.New().Interface().(*dynamicpb.Message)}
	} else {
		obj = mi.MessageType.New().Interface().(IObject)
	}
	return obj, nil
}

func (r *TypeRegistry) CreateObject(spec interface{}) (IObject, error) {
	var (
		mi *MessageInfo
		ok bool
	)

	r.mux.Lock()
	defer r.mux.Unlock()

	switch s := spec.(type) {
	case string:
		mi, ok = r.fromFullname[s]
		if !ok {
			return nil, fmt.Errorf("message with fullname [%s] not found", s)
		}
	case TypeKey:
		mi, ok = r.fromTypeKey[s]
		if !ok {
			return nil, fmt.Errorf("message with TypeKey [%X] not found", s)
		}
	default:
		return nil, fmt.Errorf("unknown type spec: [%v]", spec)
	}
	var obj IObject
	if mi.IsDynamic {
		obj = &DynObject{Message: mi.MessageType.New().Interface().(*dynamicpb.Message)}
	} else {
		obj = mi.MessageType.New().Interface().(IObject)
	}

	// Create a Metadata field
	fd := obj.ProtoReflect().Descriptor().Fields().ByName("metadata")
	newUuid, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}
	md := &Metadata{
		Type:      mi.TypeKey[:],
		Uuid:      newUuid,
		Labels:    []string{},
		CreatedAt: timestamppb.Now()}
	mdVal := protoreflect.ValueOfMessage(md.ProtoReflect())
	obj.ProtoReflect().Set(fd, mdVal)
	return obj, nil
}

type DynObject struct {
	*dynamicpb.Message
}

func (o *DynObject) GetMetadata() *Metadata {
	m := &Metadata{}
	fd := o.Descriptor().Fields().ByName("metadata")
	metaMd := fd.Message()
	if metaMd == nil {
		return nil
	}
	metaVal := o.Get(fd).Message()
	for i := 0; i < metaMd.Fields().Len(); i++ {
		switch fd := metaMd.Fields().Get(i); fd.Number() {
		case 1: // "type"
			m.Type = metaVal.Get(fd).Bytes()
		case 2: // "uuid"
			m.Uuid = metaVal.Get(fd).Bytes()
		case 3: // "labels"
			m.Labels = []string{}
			list := metaVal.Get(fd).List()
			for j := 0; j < list.Len(); j++ {
				m.Labels = append(m.Labels, list.Get(j).String())
			}
		case 4: // "description"
			m.Description = metaVal.Get(fd).String()
		case 5: // "created_at"
			ts := &timestamppb.Timestamp{}
			tsMsg := metaVal.Get(fd).Message()
			secFd := tsMsg.Descriptor().Fields().ByNumber(1)
			nanosFd := tsMsg.Descriptor().Fields().ByNumber(2)
			ts.Seconds = tsMsg.Get(secFd).Int()
			ts.Nanos = int32(tsMsg.Get(nanosFd).Int())
			m.CreatedAt = ts
		case 6: // "updated_at"
			ts := &timestamppb.Timestamp{}
			tsMsg := metaVal.Get(fd).Message()
			secFd := tsMsg.Descriptor().Fields().ByNumber(1)
			nanosFd := tsMsg.Descriptor().Fields().ByNumber(2)
			ts.Seconds = tsMsg.Get(secFd).Int()
			ts.Nanos = int32(tsMsg.Get(nanosFd).Int())
			m.CreatedAt = ts
		default:
			log.Printf("unknown field in Metadata: %s", fd.Name())
		}
	}

	return m
}

func (r *TypeRegistry) StoreSchema() error {
	r.mux.Lock()
	defer r.mux.Unlock()
	return StoreSchema(*r.files)
}

func (r *TypeRegistry) LoadSchema() (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	fds, err := LoadSchema()
	if err != nil {
		return err
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return err
	}

	r.files = files
	return r.refresh()
}

func (r *TypeRegistry) GetFileDescriptorSet() *descriptorpb.FileDescriptorSet {
	fileSet := &descriptorpb.FileDescriptorSet{}
	r.files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		fileSet.File = append(fileSet.File, protodesc.ToFileDescriptorProto(fd))
		return true
	})
	return fileSet
}

func (r *TypeRegistry) GetTypeNames() map[string][]string {
	res := map[string][]string{}
	for _, v := range r.fromTypeKey {
		res[v.Fullname] = v.Aliases
	}
	return res
}

func (r *TypeRegistry) UseFileDescriptorSet(fds *descriptorpb.FileDescriptorSet) (err error) {
	r.files = nil
	r.files, err = protodesc.NewFiles(fds)
	r.clear()
	r.refresh()
	return
}

func (r *TypeRegistry) GetTypeKeyFromToA(toa string) (TypeKey, error) {
	for k, v := range r.fromTypeKey {
		if v.Fullname == toa {
			return k, nil
		}
		for _, alias := range v.Aliases {
			if alias == toa {
				return k, nil
			}
		}
	}
	return TypeKey{}, ErrNotFound
}

func (r *TypeRegistry) Unmarshal(key []byte, value []byte) (IObject, error) {
	obj, err := r.CreateEmptyObject(TypeKey(key))
	if err != nil {
		return nil, err
	}
	if err = proto.Unmarshal(value, obj); err != nil {
		return nil, err
	}
	return obj, err
}

func (r *TypeRegistry) GetMessageInfo(tk TypeKey) (MessageInfo, error) {
	mi, ok := r.fromTypeKey[tk]
	if !ok {
		return MessageInfo{}, ErrNotFound
	}
	return *mi, nil
}
