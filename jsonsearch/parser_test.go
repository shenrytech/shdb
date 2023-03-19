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

package jsonsearch

import (
	"log"
	"strings"
	"testing"
)

const testData = `
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Compile protobuf",
            "type": "shell",
            "args": [
                "${workspaceFolder}"
            ],
            "problemMatcher": [
                "$tsc"
            ],
            "windows": {
                "command": "${workspaceFolder}/scripts/build_protobuf.ps1"
            },
            "linux": {
                "command": "${workspaceFolder}/scripts/build_protobuf.sh"
            },
            "presentation": {
                "reveal": "silent"
            },
            "group": "build"
        }
    ]
}`

func TestJsonParser(t *testing.T) {
	p := NewParser([]byte(testData), func(s string) bool {
		return strings.Contains(s, "workspaceFolder")
	})
	err := p.Parse("")
	if err != nil {
		log.Fatalf("parser failed %v", err)
		t.Fail()
	}

	expected := []string{"/tasks/@0/args/@0", "/tasks/@0/windows/command", "/tasks/@0/linux/command"}

	for k := range expected {
		if expected[k] != p.FieldPaths[k] {
			t.Fail()
		}
	}

}
