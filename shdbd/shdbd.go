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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/shenrytech/shdb"
	"google.golang.org/grpc"
)

var uuids = []string{
	"025991b3-67ed-435b-8fad-bf42d757ecd1",
	"078038e6-a5f4-432a-a129-690689dcc3ca",
	"0c8b6424-a829-4af5-9a89-8fc02aa5de55",
	"295c81a9-66c1-4c00-b7d0-77f62927d7c7",
	"2eb4f2e7-ad0b-4416-92f1-f7007b4e5738",
	"315beabd-21a7-491a-8d8d-ad754faddbcc",
	"590cc0bd-2222-4de9-aece-75c340900ea5",
	"5f6dc7f9-5f82-42dd-b33e-1d86e7f82944",
	"645caaae-b2c1-4a0c-8558-a2953c7f0065",
	"6f5d8f6d-6aa4-4b6d-a5dd-db323bdcde96",
	"71a13b7b-15bd-4dfa-bb58-26c203c73f19",
	"80a896b4-2f23-4850-aa73-5a74628ea3e2",
	"92d0bac2-bb0f-4a11-843b-4a586986e361",
	"9eaaad83-560d-41e4-9139-d3507fb51d13",
	"abec119d-d089-4d1e-aae6-7f5df16f78d2",
	"c169dfef-211a-4477-832a-18555f0a9bff",
	"d6ed960f-8cef-4583-bb3e-0daee1c765b0",
	"e1f49b04-8f15-4241-a420-b13340779efd",
	"f4a78ea5-3149-4755-9d02-a37062463290",
	"fe4ceb7c-b67d-4893-9cc8-fd287dd73bf8",
}

var uuidbs [][]byte

func uuidb(idx int) []byte {
	if uuidbs == nil {
		uuidbs = [][]byte{}
		for _, v := range uuids {
			data, err := uuid.MustParse(v).MarshalBinary()
			if err != nil {
				panic(err)
			}
			uuidbs = append(uuidbs, data)
		}
	}
	return uuidbs[idx]
}

func loadtd() {
	TTObject := shdb.TypeKeyOf("shdb.TObject")
	if err := shdb.DeleteAll(TTObject); err != nil {
		panic(err)
	}
	for i := 0; i < len(uuids); i++ {
		obj, err := shdb.New[*shdb.TObject](TTObject)
		if err != nil {
			panic(err)
		}
		obj.Metadata.Uuid = uuidb(i)
		obj.MyInt = uint64(i)
		if err := shdb.Put(obj); err != nil {
			panic(err)
		}
	}

}

func main() {
	serverPort := flag.Int("grpc-port", 3335, "api server port to listen on")
	dbFile := flag.String("dbfile", "/tmp/shdb.db", "database file")
	loadTestData := flag.Bool("load-test-data", false, "load test data")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *serverPort))
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	shdb.Init(*dbFile)
	if *loadTestData {
		loadtd()
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	shdb.NewServer(context.Background(), grpcServer, shdb.NewTypeRegistry())
	grpcServer.Serve(listener)
}
