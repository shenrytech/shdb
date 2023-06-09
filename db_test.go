// Copyright 2023 Shenry Tech AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package shdb

import (
	"context"
	"errors"
	"os"
	"path"
	"sort"
	"testing"

	"google.golang.org/protobuf/proto"
)

var (
	TObj = TypeKeyOf("shdb.v1.TObject")
)

func CreateTestDb() string {
	tmpDir, err := os.MkdirTemp("", "shdb_test")
	if err != nil {
		panic(err)
	}
	Init(path.Join(tmpDir, "test.db"))
	return tmpDir
}

func CloseTestDb(tmpDir string) {
	Close()
	os.RemoveAll(tmpDir)
}

func GenerateTestData(count int) ([]*TObject, string) {
	tmpDir, err := os.MkdirTemp("", "shdb_test")
	if err != nil {
		panic(err)
	}
	Init(path.Join(tmpDir, "test.db"))

	list := []*TObject{}
	for k := 0; k < count; k++ {
		tObj := MustNew[*TObject](TObj)
		tObj.GetMetadata().Description = "Staffan Olsson was here"
		tObj.MyInt = uint64(k)
		list = append(list, tObj)
	}
	if err := Put(list...); err != nil {
		panic(err)
	}
	return list, tmpDir
}

func RemoveTestData(tmpDir string) {
	Close()
	os.RemoveAll(tmpDir)
}

func CompareSame(a, b []*TObject) bool {
	if len(a) != len(b) {
		return false
	}
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].MyInt < a[j].MyInt
	})
	sort.SliceStable(b, func(i, j int) bool {
		return b[i].MyInt < b[j].MyInt
	})
	for k := range a {
		if !proto.Equal(a[k], b[k]) {
			return false
		}
	}
	return true
}

func TestDB(t *testing.T) {
	dbFile := path.Join(os.TempDir(), "db_test.db")
	Init(dbFile)
	defer func() {
		Close()
		os.Remove(dbFile)
	}()

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
	list, testDir := GenerateTestData(1000)
	defer RemoveTestData(testDir)
	ctx := context.Background()
	list2, nextToken, err := List[*TObject](ctx, TObj, 400000, "")
	if err != nil {
		t.Fail()
	}
	if nextToken != "" {
		t.Fail()
	}

	if !CompareSame(list, list2) {
		t.Fail()
	}

	var (
		nextPageToken string = ""
		list3         []*TObject
	)
	list4 := []*TObject{}
	for {
		list3, nextPageToken, err = List[*TObject](ctx, TObj, 1000, nextPageToken)
		if err != nil {
			t.Fail()
		}
		list4 = append(list4, list3...)
		if nextPageToken == "" {
			break
		}
	}

	if !CompareSame(list, list4) {
		t.Fail()
	}
}

func BenchmarkListLong(b *testing.B) {

	list, testDir := GenerateTestData(1000)
	defer RemoveTestData(testDir)

	var (
		nextPageToken string = ""
		list3         []*TObject
		err           error
	)
	ctx := context.Background()
	list4 := []*TObject{}
	for {
		list3, nextPageToken, err = List[*TObject](ctx, TObj, 400000, nextPageToken)
		if err != nil {
			b.Fail()
		}
		list4 = append(list4, list3...)
		if nextPageToken == "" {
			break
		}
	}

	if !CompareSame(list, list4) {
		b.Fail()
	}
}
