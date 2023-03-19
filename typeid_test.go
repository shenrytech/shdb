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
	"bytes"
	"testing"

	"github.com/google/uuid"
)

var (
	idBytes = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	typeKey = TypeKey{22, 23, 24, 25}
)

func TestTypeId(t *testing.T) {
	id, _ := uuid.FromBytes(idBytes)

	dbk := NewTypeId(typeKey, idBytes)

	if !bytes.Equal(dbk.data[:], []byte{22, 23, 24, 25, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}) {
		t.Fail()

	}

	if dbk.Uuid() != id {
		t.Fail()
	}

	dbkString := dbk.String()
	if dbkString != "%16%17%18%19%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10" {
		t.Fail()
	}

	dbk2, err := TypeIdFromString((dbkString))
	if err != nil {
		t.Fail()
	}

	if !dbk2.Equal(dbk) {
		t.Fail()
	}

	dbkt := dbk.TypeKey()
	if !bytes.Equal(dbkt[:], []byte{22, 23, 24, 25}) {
		t.Fail()
	}
}
