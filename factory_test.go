package shdb

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestGet(t *testing.T) {
	var TObj = TypeKey{0, 0, 0, 1}
	Register(&TObject{Metadata: &Metadata{Type: TObj[:]}})
	a := MustNew[*TObject](TObj)
	data, _ := Marshal(a)

	b, _ := Unmarshal[*TObject](data[0])
	if !proto.Equal(a, b) {
		t.Fail()
	}
}
