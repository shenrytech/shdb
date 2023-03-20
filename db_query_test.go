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
	"context"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	count := 10
	pageSize := 10
	list, testDir := GenerateTestData(count)
	defer RemoveTestData(testDir)

	var (
		nextPageToken string = ""
		list3         []*TObject
		err           error
	)
	ctx := context.Background()

	identityFn := func(obj *TObject) (bool, error) {
		return true, nil
	}
	list4 := []*TObject{}
	for {
		list3, nextPageToken, err = Query(ctx, TObj, identityFn, int32(count/pageSize), nextPageToken)
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

func TestQueryPageToken(t *testing.T) {
	count := 100
	pageSize := 10
	list, testDir := GenerateTestData(count)
	defer RemoveTestData(testDir)

	var (
		nextPageToken string = ""
		list3         []*TObject
		err           error
	)
	ctx := context.Background()

	identityFn := func(obj *TObject) (bool, error) {
		return true, nil
	}
	list4 := []*TObject{}
	for {
		list3, nextPageToken, err = Query(ctx, TObj, identityFn, int32(count/pageSize), nextPageToken)
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

func TestQueryFilter(t *testing.T) {
	count := 100
	pageSize := 10
	list, testDir := GenerateTestData(count)
	defer RemoveTestData(testDir)

	var (
		nextPageToken string = ""
		list3         []*TObject
		err           error
	)
	ctx := context.Background()

	evenFn := func(obj *TObject) (bool, error) {
		return obj.MyInt%2 == 0, nil
	}
	list4 := []*TObject{}
	for {
		list3, nextPageToken, err = Query(ctx, TObj, evenFn, int32(count/pageSize), nextPageToken)
		if err != nil {
			t.Fail()
		}
		list4 = append(list4, list3...)
		if nextPageToken == "" {
			break
		}
	}

	evenList := []*TObject{}
	for _, v := range list {
		if v.MyInt%2 == 0 {
			evenList = append(evenList, v)
		}
	}

	if !CompareSame(evenList, list4) {
		t.Fail()
	}
}

func ExampleQuery() {
	Init(path.Join(os.TempDir(), "example_query.db"))
	Register(&TObject{
		Metadata: &Metadata{Type: TObj[:]},
		MyField:  "The flying duck is flying low"})

	count := 100
	pageSize := 10

	list := []*TObject{}
	for k := 0; k < count; k++ {
		tObj := MustNew[*TObject](TObj)
		tObj.MyInt = uint64(k)
		list = append(list, tObj)
	}
	if err := Put(list...); err != nil {
		panic(err)
	}

	defer func() {
		Close()
		os.Remove(path.Join(os.TempDir(), "example_query.db"))
	}()

	var (
		nextPageToken string = ""
		partList      []*TObject
		err           error
	)
	ctx := context.Background()

	selectorFn := func(obj *TObject) (bool, error) {
		return strings.Contains(obj.MyField, "flying"), nil
	}
	collected := []*TObject{}
	for {
		partList, nextPageToken, err = Query(ctx, TObj, selectorFn, int32(count/pageSize), nextPageToken)
		if err != nil {
			panic(err)
		}
		collected = append(collected, partList...)
		if nextPageToken == "" {
			break
		}
	}
	for idx, obj := range collected {
		log.Printf("%d: [%v]\n", idx, obj)
	}
}
