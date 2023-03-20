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
	"fmt"
	"log"
	"testing"
)

func TestWatch(t *testing.T) {
	count := 10
	ch := make(chan *EventInfo, 1)
	watchId, err := WatchType("", ch, TObj)
	if err != nil {
		t.FailNow()
	}
	go func() {
		defer UnwatchType(watchId, TObj)
		for i := 0; i < 10; i++ {
			ev := <-ch
			if ev.Kind != EventCreated {
				t.Fail()
			}
			tobj := ev.Object.(*TObject)
			if tobj.MyInt != uint64(i) {
				t.Fail()
			}
		}
	}()
	_, testDir := GenerateTestData(count)
	defer RemoveTestData(testDir)
	if err = RemoveWatcher(watchId); err != nil {
		t.Fail()
	}
}

func TestWatch2(t *testing.T) {
	count := 10
	ch := make(chan *EventInfo, 1)
	watchId, err := WatchType("", ch, TObj)
	if err != nil {
		t.FailNow()
	}
	chDone := make(chan error)
	go func() {
		// Wait for created objects
		for i := 0; i < count; i++ {
			ev := <-ch
			if ev.Kind != EventCreated {
				chDone <- fmt.Errorf("ev.Kind != EventCreated")
				return
			}
			tobj := ev.Object.(*TObject)
			if tobj.MyInt != uint64(i) {
				chDone <- fmt.Errorf("tobj.MyInt != uint64(i)")
				return
			}
		}

		// Wait for updated objects
		for i := 0; i < count; i++ {
			ev := <-ch
			if ev.Kind != EventUpdated {
				chDone <- fmt.Errorf("ev.Kind != EventUpdate")
				return
			}
			tobj := ev.Object.(*TObject)
			tprev := ev.Previous.(*TObject)
			if tobj.MyInt != tprev.MyInt+1 {
				chDone <- fmt.Errorf("tobj.MyInt != tprev.MyInt-1")
				return
			}
		}

		// Wait for deleted objects
		for i := 0; i < count; i++ {
			ev := <-ch
			if ev.Kind != EventDeleted {
				chDone <- fmt.Errorf("ev.Kind != EventDeleted")
				return
			}
			tobj := ev.Object.(*TObject)
			if tobj.MyInt != uint64(i+1) {
				chDone <- fmt.Errorf("tobj.MyInt != uint64(i+1)")
				return
			}
		}
		chDone <- nil
	}()

	_, testDir := GenerateTestData(count)
	defer RemoveTestData(testDir)

	// Update
	for i := count - 1; i >= 0; i-- {
		obj, err := GetFirst(TObj, func(obj *TObject) bool {
			return obj.MyInt == uint64(i)
		})
		if err != nil {
			t.FailNow()
		}
		_, err = Update(obj.Metadata.TypeId(), func(obj *TObject) (*TObject, error) {
			obj.MyInt++
			return obj, nil
		})
		if err != nil {
			t.FailNow()
		}
	}

	// Delete
	for i := 1; i <= count; i++ {
		obj, err := GetFirst(TObj, func(obj *TObject) bool {
			return obj.MyInt == uint64(i)
		})
		if err != nil {
			t.FailNow()
		}
		_, err = Delete[*TObject](obj.Metadata.TypeId())
		if err != nil {
			t.FailNow()
		}
	}
	err = <-chDone
	if err != nil {
		log.Printf("check failed %v\n", err)
	}
	UnwatchType(watchId, TObj)
}
