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
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// Events reflecting life-cycle changes to an object
const (
	InvalidWatchEvent = 0
	EventCreated      = 1
	EventUpdated      = 2
	EventDeleted      = 3
)

// EventInfo contains the information about a specific life cycle event
type EventInfo struct {
	Kind     int
	Tid      TypeId
	Object   IObject
	Previous IObject
}

type watchInstance struct {
	Types   []TypeKey
	TypeIds []TypeId
	Ch      chan *EventInfo
}

type watchCtrlRsp struct {
	watcherId string
	err       error
}

type watchCtrlReq struct {
	watcherId  string
	addTypes   []TypeKey
	addTypeIds []TypeId
	rmTypes    []TypeKey
	rmTypeIds  []TypeId
	rmWatcher  bool
	evCh       chan *EventInfo
	rsp        chan watchCtrlRsp
}

var (
	cmdCh chan watchCtrlReq
	evCh  chan *EventInfo
)

func watchRun(ctx context.Context) {

	watchInstances := map[string]*watchInstance{}

	handleCmd := func(cmd watchCtrlReq) {
		rsp := watchCtrlRsp{watcherId: cmd.watcherId}
		if cmd.addTypeIds != nil || cmd.addTypes != nil {
			if cmd.watcherId == "" {
				rsp.watcherId = uuid.NewString()
				watchInstances[rsp.watcherId] = &watchInstance{Ch: cmd.evCh}
			}
		}
		if cmd.rmWatcher {
			close(watchInstances[cmd.watcherId].Ch)
			delete(watchInstances, rsp.watcherId)
			cmd.rsp <- rsp
			return
		}
		if cmd.addTypeIds != nil {
			watchInstances[rsp.watcherId].TypeIds = append(watchInstances[rsp.watcherId].TypeIds, cmd.addTypeIds...)
		}
		if cmd.addTypes != nil {
			watchInstances[rsp.watcherId].Types = append(watchInstances[rsp.watcherId].Types, cmd.addTypes...)
		}
		if cmd.rmTypeIds != nil {
			res := []TypeId{}
		tidLoop:
			for _, v := range watchInstances[rsp.watcherId].TypeIds {
				for _, rv := range cmd.rmTypeIds {
					if v == rv {
						continue tidLoop
					}
				}
				res = append(res, v)
			}
			watchInstances[rsp.watcherId].TypeIds = res
		}
		if cmd.rmTypes != nil {
			res := []TypeKey{}
		tLoop:
			for _, v := range watchInstances[rsp.watcherId].Types {
				for _, rv := range cmd.rmTypes {
					if v == rv {
						continue tLoop
					}
				}
				res = append(res, v)
			}
			watchInstances[rsp.watcherId].Types = res
		}
		cmd.rsp <- rsp
	}

	handleEvent := func(ev *EventInfo) {
		typ := ev.Tid.TypeKey()
		for _, v := range watchInstances {
			for _, wv := range v.Types {
				if wv == typ {
					v.Ch <- ev
				}
			}
			for _, wv := range v.TypeIds {
				if wv == ev.Tid {
					v.Ch <- ev
				}
			}
		}
	}

	closeAll := func() {
		for _, v := range watchInstances {
			close(v.Ch)
		}
	}

	defer closeAll()

	for {
		select {
		case cmd := <-cmdCh:
			handleCmd(cmd)
		case <-ctx.Done():
			return
		case ev := <-evCh:
			handleEvent(ev)
		}
	}
}

func notifyCreate(obj IObject) {
	m := obj.(proto.Message)
	ev := &EventInfo{
		Kind:     EventCreated,
		Object:   proto.Clone(m).(IObject),
		Previous: nil,
		Tid:      obj.GetMetadata().TypeId(),
	}
	evCh <- ev
}

func notifyUpdate(obj, prev IObject) {
	mObj := obj.(proto.Message)
	mPrev := prev.(proto.Message)
	ev := &EventInfo{
		Kind:     EventUpdated,
		Object:   proto.Clone(mObj).(IObject),
		Previous: proto.Clone(mPrev).(IObject),
		Tid:      obj.GetMetadata().TypeId(),
	}
	evCh <- ev
}

func notifyDelete(obj IObject) {
	mObj := obj.(proto.Message)
	ev := &EventInfo{
		Kind:     EventDeleted,
		Object:   proto.Clone(mObj).(IObject),
		Previous: nil,
		Tid:      obj.GetMetadata().TypeId(),
	}
	evCh <- ev
}

// WatchType creates or updates a watcher by adding watches to new TypeKeys
// If the provided watcherId is the empty string, a new watcher is created and the
// eventCh must be specified. If watcherId is non-empty, then the eventCh can be set to nil
// The watcherId is returned.
func WatchType(watcherId string, eventCh chan *EventInfo, typeKeys ...TypeKey) (string, error) {
	req := watchCtrlReq{
		watcherId: watcherId,
		addTypes:  typeKeys,
		evCh:      eventCh,
		rsp:       make(chan watchCtrlRsp),
	}
	cmdCh <- req
	rsp := <-req.rsp
	close(req.rsp)
	return rsp.watcherId, rsp.err
}

// UnwatchType removes a list of TypeKeys from a watcher
func UnwatchType(watcherId string, typeKeys ...TypeKey) error {
	if watcherId == "" {
		return ErrSessionInvalid
	}
	req := watchCtrlReq{
		watcherId: watcherId,
		rmTypes:   typeKeys,
		rsp:       make(chan watchCtrlRsp),
	}
	cmdCh <- req
	rsp := <-req.rsp
	close(req.rsp)
	return rsp.err
}

// WatchType creates or updates a watcher by adding watches to new TypeIds
// If the provided watcherId is the empty string, a new watcher is created and the
// eventCh must be specified. If watcherId is non-empty, then the eventCh can be set to nil
// The watcherId is returned.
func WatchTypeId(watcherId string, eventCh chan *EventInfo, tids ...TypeId) (string, error) {
	req := watchCtrlReq{
		watcherId:  watcherId,
		addTypeIds: tids,
		evCh:       eventCh,
		rsp:        make(chan watchCtrlRsp),
	}
	cmdCh <- req
	rsp := <-req.rsp
	close(req.rsp)
	return rsp.watcherId, rsp.err
}

// UnwatchTypeId removes a list of TypeIds from a watcher
func UnwatchTypeId(watcherId string, tids ...TypeId) error {
	if watcherId == "" {
		return ErrSessionInvalid
	}
	req := watchCtrlReq{
		watcherId: watcherId,
		rmTypeIds: tids,
		rsp:       make(chan watchCtrlRsp),
	}
	cmdCh <- req
	rsp := <-req.rsp
	close(req.rsp)
	return rsp.err
}

// RemoveWatcher closes the eventCh for the watcher and removes
// the watcher.
func RemoveWatcher(watcherId string) error {
	if watcherId == "" {
		return ErrSessionInvalid
	}
	req := watchCtrlReq{
		watcherId: watcherId,
		rmWatcher: true,
		rsp:       make(chan watchCtrlRsp),
	}
	cmdCh <- req
	rsp := <-req.rsp
	close(req.rsp)
	return rsp.err
}

func init() {
	evCh = make(chan *EventInfo, 1)
	cmdCh = make(chan watchCtrlReq, 1)
	go watchRun(context.Background())
	// fmt.Println("init")
}
