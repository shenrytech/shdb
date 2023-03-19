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
