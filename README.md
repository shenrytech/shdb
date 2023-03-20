# Shdb

A simple database for protobuf based objects

## Installing

```sh
$ go get github.com/shenrytech/shdb
```


## Example

```go
package main

func main() {
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
		partList, nextPageToken, err = Query(ctx, 
            TObj, selectorFn, int32(count/pageSize), nextPageToken)
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
```

## License

Shdb is released under the Apache 2.0 license. See LICENSE.txt

