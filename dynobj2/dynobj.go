package dynobj2

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"gopkg.in/yaml.v3"
)

func MergeStrings(a, b []string) []string {
	merger := func(list []string, str string) []string {
		for _, v := range a {
			if v == str {
				return list
			}
		}
		return append(list, str)
	}
	for _, v := range b {
		a = merger(a, v)
	}
	return a
}

type TypeInfo struct {
	FullName    string
	Aliases     []string
	ObjInstance proto.Message
	Descriptor  descriptorpb.DescriptorProto
}

type TReg struct {
	fileSet     map[string]*descriptorpb.FileDescriptorProto
	filePackMap map[string]string
	files       *protoregistry.Files
}

func NewTReg() (*TReg, error) {
	r := &TReg{
		fileSet:     make(map[string]*descriptorpb.FileDescriptorProto),
		filePackMap: map[string]string{},
		files:       nil,
	}
	err := r.LoadGlobals()
	return r, err
}

func (r *TReg) dumpDependencies() {
	for pack, fd := range r.fileSet {
		fmt.Printf("%s [%s]\n", fd.GetName(), pack)
		for _, v := range fd.Dependency {
			fmt.Printf(" -->%s\n", v)
		}
	}
}

func (r *TReg) refreshFiles() (err error) {
	// Update the import path
	for _, fd := range r.fileSet {
		deps := map[string]struct{}{}
		for _, v := range fd.Dependency {
			if !strings.HasPrefix(v, "shdb_dyn") {
				v = fmt.Sprintf("shdb_dyn/%s.proto", r.filePackMap[v])
			}
			if v != fd.GetName() {
				deps[v] = struct{}{}
			}
		}
		fd.Dependency = []string{}
		for k := range deps {
			fd.Dependency = append(fd.Dependency, k)
		}
	}
	r.dumpDependencies()
	fileSet := &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{},
	}
	for _, v := range r.fileSet {
		fileSet.File = append(fileSet.File, v)
	}
	r.files, err = protodesc.NewFiles(fileSet)
	return err
}

func (r *TReg) LoadGlobals() error {
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {

		// Since gRPC is using old github.com/golang/protobuf, we need to
		// filter out the xtra github protobuf files that is added by that package
		if strings.HasPrefix(fd.Path(), "github.com/golang/protobuf/ptypes") {
			return true
		}
		r.filePackMap[fd.Path()] = string(fd.Package())
		pack := string(fd.Package())
		var (
			fdp *descriptorpb.FileDescriptorProto
			ok  bool
		)
		fdp, ok = r.fileSet[pack]
		log.Printf("loading %s\n", fd.Path())
		if !ok {
			log.Printf("New proto package %s\n", pack)
			fName := fmt.Sprintf("shdb_dyn/%s.proto", pack)
			fdp = &descriptorpb.FileDescriptorProto{
				Package: &pack,
				Name:    &fName,
			}
			r.fileSet[pack] = fdp
		} else {
			log.Printf("Merging proto package %s\n", pack)
		}
		r.MergeInto(fdp, protodesc.ToFileDescriptorProto(fd))
		return true
	})
	r.dumpDependencies()
	return r.refreshFiles()
}

func (r *TReg) MergeInto(fd, other *descriptorpb.FileDescriptorProto) error {
	if other == nil {
		return nil
	}
	if fd.Dependency == nil {
		fd.Dependency = []string{}
	}
	if fd.PublicDependency == nil {
		fd.PublicDependency = []int32{}
	}
	if fd.MessageType == nil {
		fd.MessageType = []*descriptorpb.DescriptorProto{}
	}
	if fd.Options == nil {
		fd.Options = &descriptorpb.FileOptions{}
	}
	if fd.EnumType == nil {
		fd.EnumType = []*descriptorpb.EnumDescriptorProto{}
	}
	if fd.Service == nil {
		fd.Service = []*descriptorpb.ServiceDescriptorProto{}
	}

	// Merge dependencies and PublicDependencies

	// Collect all package dependencies of 'others'
depsLoop:
	for _, otherDep := range other.Dependency {
		for k, thisDep := range fd.Dependency {
			if otherDep == thisDep {
				// Check if this is a public dependency
			pubDepLoop:
				for _, otherDepIdx := range other.PublicDependency {
					if int32(k) == otherDepIdx {
						// It was. Now see if it's already in fd.PublicDependency
						for _, thisDepIdx := range fd.PublicDependency {
							if thisDepIdx == otherDepIdx {
								break pubDepLoop
							}
						}
						fd.PublicDependency = append(fd.PublicDependency, int32(k))
						break pubDepLoop
					}
				}
				continue depsLoop
			}
		}
		fd.Dependency = append(fd.Dependency, otherDep)
	}

	// Merge Messages
messageLoop:
	for _, otherMsg := range other.MessageType {
		for _, thisMsg := range fd.MessageType {
			if *otherMsg.Name == *thisMsg.Name {
				continue messageLoop
			}
		}
		newMsg := proto.Clone(otherMsg).(*descriptorpb.DescriptorProto)
		fd.MessageType = append(fd.MessageType, newMsg)
	}

	// Merge Enums
enumLoop:
	for _, otherItem := range other.EnumType {
		for _, thisItem := range fd.EnumType {
			if *otherItem.Name == *thisItem.Name {
				continue enumLoop
			}
		}
		newItem := proto.Clone(otherItem).(*descriptorpb.EnumDescriptorProto)
		fd.EnumType = append(fd.EnumType, newItem)
	}

	// Merge Services
srvcLoop:
	for _, otherItem := range other.Service {
		for _, thisItem := range fd.Service {
			if *otherItem.Name == *thisItem.Name {
				continue srvcLoop
			}
		}
		newItem := proto.Clone(otherItem).(*descriptorpb.ServiceDescriptorProto)
		fd.Service = append(fd.Service, newItem)
	}

	// Merge Extensions
extLoop:
	for _, otherItem := range other.Extension {
		for _, thisItem := range fd.Extension {
			if *otherItem.Name == *thisItem.Name {
				continue extLoop
			}
		}
		newItem := proto.Clone(otherItem).(*descriptorpb.FieldDescriptorProto)
		fd.Extension = append(fd.Extension, newItem)
	}
	log.Println(fd.Dependency)
	return nil
}

func (r *TReg) AddMessageFromYaml(pack string, data []byte) error {
	var (
		fdp *descriptorpb.FileDescriptorProto
		ok  bool
	)
	fdp, ok = r.fileSet[pack]
	if !ok {
		fName := fmt.Sprintf("shdb_dyn/%s.proto", pack)
		fdp = &descriptorpb.FileDescriptorProto{
			Package: &pack,
			Name:    &fName,
		}
		r.fileSet[pack] = fdp
	}

	msgDescr := &descriptorpb.DescriptorProto{}
	err := yaml.Unmarshal(data, msgDescr)
	if err != nil {
		return err
	}
	for k, v := range fdp.MessageType {
		if v.GetName() == msgDescr.GetName() {
			log.Printf("replacing %s", v.GetName())
			fdp.MessageType = append(fdp.MessageType[:k], fdp.MessageType[k+1:]...)
			break
		}
		log.Printf("creating %s", v.GetName())
	}
	fdp.MessageType = append(fdp.MessageType, msgDescr)
	return r.refreshFiles()
}

func (r *TReg) CreateMessage(fullname string) (proto.Message, error) {
	if r.files == nil {
		return nil, fmt.Errorf("files not initialized")
	}
	descr, err := r.files.FindDescriptorByName(protoreflect.FullName(fullname))
	if err != nil {
		return nil, err
	}
	msgDesc, ok := descr.(protoreflect.MessageDescriptor)
	if !ok {
		return nil, fmt.Errorf("%s is not a message", fullname)
	}
	return dynamicpb.NewMessage(msgDesc), nil
}

func GenerateFileSet() (*descriptorpb.FileDescriptorSet, error) {
	set := &descriptorpb.FileDescriptorSet{
		File: make([]*descriptorpb.FileDescriptorProto, 0),
	}

	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		set.File = append(set.File, protodesc.ToFileDescriptorProto(fd))
		return true
	})

	data, err := yaml.Marshal(set)
	if err != nil {
		return nil, err
	}
	fmt.Print(string(data))
	return set, nil
}

type ListAllMessagesOutput struct {
	Fullname   string
	ParentFile string
}

func (r *TReg) ListAllMessages() []ListAllMessagesOutput {
	res := []ListAllMessagesOutput{}

	r.files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Messages().Len(); i++ {
			msg := fd.Messages().Get(i)
			res = append(res, ListAllMessagesOutput{
				Fullname:   string(msg.FullName()),
				ParentFile: msg.ParentFile().Path(),
			})
		}
		return true
	})
	return res
}
