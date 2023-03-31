package dynobj2

type MessageInfo struct {
	FullName      string
	Aliases       []string
	DescriptorRaw []byte
}

// A MessageFile contains all information needed
// to generate a FileDescriptor containing all thats
// needed to generate protobuf messages compatible with
// the ones created by the compiler.
type MessageFile struct {
	Path    string
	Package string
	Imports []string
}
