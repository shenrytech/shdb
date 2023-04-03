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
	"errors"
	"log"

	"go.etcd.io/bbolt"
)

var (
	bucket_obj    = []byte("obj")
	bucket_schema = []byte("schema")
	db            *bbolt.DB
	typeRegistry  *TypeRegistry
)

// Init initializes the backing database.
func Init(dbFile string) {
	var err error
	db, err = bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists(bucket_obj)
		tx.CreateBucketIfNotExists(bucket_schema)
		return nil
	})
	typeRegistry = NewTypeRegistry()
	err = typeRegistry.LoadSchema()
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			err = typeRegistry.refresh()
			log.Println("loaded schema from runtime")
			if err != nil {
				panic(err)
			}
		}
	} else {
		log.Println("loaded schema from database")
	}
}

// Close the backing database
func Close() error {
	return db.Close()
}
