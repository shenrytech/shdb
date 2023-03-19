package shdb

import (
	"context"
	"testing"
)

func TestQuery(t *testing.T) {
	list := GenerateTestData(10)
	defer RemoveTestData()

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
		list3, nextPageToken, err = Query(ctx, TObj, identityFn, 400000, nextPageToken)
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
	list := GenerateTestData(10000)
	defer RemoveTestData()

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
		list3, nextPageToken, err = Query(ctx, TObj, identityFn, 1, nextPageToken)
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
	list := GenerateTestData(10000)
	defer RemoveTestData()

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
		list3, nextPageToken, err = Query(ctx, TObj, evenFn, 1, nextPageToken)
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
