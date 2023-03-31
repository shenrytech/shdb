package dynobj2

import (
	"fmt"
	"log"
	"testing"

	"github.com/shenrytech/shdb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func TestGenerateFileSet(t *testing.T) {
	_ = &shdb.TObject{}

	msgType, err := protoregistry.GlobalTypes.FindMessageByName("shdb.TObject")
	if err != nil {
		t.FailNow()
	}
	msg1 := msgType.New()
	apa := msg1.Interface().(*shdb.TObject)
	log.Println(apa)

	r, err := NewTReg()
	if err != nil {
		t.FailNow()
	}
	if len(r.fileSet) == 0 {
		t.Fail()
	}
	res := r.ListAllMessages()
	for _, r := range res {
		fmt.Printf("%s [%s]\n", r.Fullname, r.ParentFile)
	}

	msg, err := r.CreateMessage("shdb.TObject")
	if err != nil {
		t.FailNow()
	}
	log.Println(msg)
	fmt.Printf("msg.ProtoReflect().Descriptor(): %v\n", msg.ProtoReflect().Descriptor())

	tobj := &shdb.TObject{}
	fmt.Println(tobj.ProtoReflect().Descriptor())
	md := tobj.ProtoReflect().Descriptor()
	a := proto.GetExtension(md.Options(), shdb.E_ShdbOptions)
	fmt.Println(a)
}
