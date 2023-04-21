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
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestGet(t *testing.T) {
	tmpDir := CreateTestDb()
	defer CloseTestDb(tmpDir)

	TObj := TypeKeyOf("shdb.v1.TObject")
	a := MustNew[*TObject](TObj)
	data, _ := Marshal(a)

	b, _ := Unmarshal[*TObject](data[0])
	if !proto.Equal(a, b) {
		t.Fail()
	}
}
