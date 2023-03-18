package shdb

import (
	"errors"
	"os"
	"path"
	"sort"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	TObj = TypeKey{1, 2, 3, 4}
)

func TestDB(t *testing.T) {
	dbFile := path.Join(os.TempDir(), "db_test.db")
	l, _ := zap.NewDevelopment()
	Init(l, dbFile)
	defer func() {
		Close()
		os.Remove(dbFile)
	}()
	Register(&TObject{Metadata: &Metadata{Type: TObj[:]}})
	t1 := MustNew[*TObject](TObj)

	if t1.GetMetadata().CreatedAt.Seconds == 0 {
		t.Fail()
	}

	if err := Put(t1); err != nil {
		t.Fail()
	}

	t2, err := Get[*TObject](t1.Metadata.TypeId())
	if err != nil {
		t.Fail()
	}

	if !proto.Equal(t1, t2) {
		t.Fail()
	}

	prev, err := Delete[*TObject](t1.Metadata.TypeId())
	if err != nil {
		t.Fail()
	}
	if !proto.Equal(t1, prev) {
		t.Fail()
	}

	_, err = Get[*TObject](t1.Metadata.TypeId())
	if err == nil {
		t.Fail()
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fail()
	}
}

func TestList(t *testing.T) {
	dbFile := path.Join(os.TempDir(), "db_test.db")
	l, _ := zap.NewDevelopment()
	Init(l, dbFile)
	defer func() {
		Close()
		os.Remove(dbFile)
	}()
	Register(&TObject{Metadata: &Metadata{Type: TObj[:]}})

	list := []*TObject{}
	for k := 0; k < 2000; k++ {
		tObj := MustNew[*TObject](TObj)
		tObj.MyInt = uint64(k)
		list = append(list, tObj)
	}
	if err := Put(list...); err != nil {
		t.Fail()
	}

	list2, nextToken, err := List[*TObject](TObj, 400000, "")
	if err != nil {
		t.Fail()
	}
	if nextToken != "" {
		t.Fail()
	}
	sort.SliceStable(list2, func(i, j int) bool {
		return list2[i].MyInt < list2[j].MyInt
	})
	for k := range list2 {
		if !proto.Equal(list[k], list2[k]) {
			t.Fail()
		}
	}

	var (
		nextPageToken string = ""
		list3         []*TObject
	)
	list4 := []*TObject{}
	for {
		list3, nextPageToken, err = List[*TObject](TObj, 1000, nextPageToken)
		if err != nil {
			t.Fail()
		}
		list4 = append(list4, list3...)
		if nextPageToken == "" {
			break
		}
	}

	sort.SliceStable(list4, func(i, j int) bool {
		return list4[i].MyInt < list4[j].MyInt
	})
	for k := range list4 {
		if !proto.Equal(list[k], list4[k]) {
			t.Fail()
		}
	}
}

func BenchmarkListLong(b *testing.B) {
	dbFile := path.Join(os.TempDir(), "db_test.db")
	l, _ := zap.NewDevelopment()
	Init(l, dbFile)
	defer func() {
		Close()
		os.Remove(dbFile)
	}()
	Register(&TObject{Metadata: &Metadata{Type: TObj[:]}})

	list := []*TObject{}
	for k := 0; k < 2000; k++ {
		tObj := MustNew[*TObject](TObj)
		tObj.MyInt = uint64(k)
		list = append(list, tObj)
	}
	if err := Put(list...); err != nil {
		b.Fail()
	}

	var (
		nextPageToken string = ""
		list3         []*TObject
		err           error
	)
	list4 := []*TObject{}
	for {
		list3, nextPageToken, err = List[*TObject](TObj, 400000, nextPageToken)
		if err != nil {
			b.Fail()
		}
		list4 = append(list4, list3...)
		if nextPageToken == "" {
			break
		}
	}

	sort.SliceStable(list4, func(i, j int) bool {
		return list4[i].MyInt < list4[j].MyInt
	})
	for k := range list4 {
		if !proto.Equal(list[k], list4[k]) {
			b.Fail()
		}
	}
}
