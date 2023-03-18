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
