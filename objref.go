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

// Marshal an ObjRef to a byte slice
func (r *ObjRef) Marshal() []byte {
	res := []byte{}
	res = append(res, r.Type...)
	res = append(res, r.Uuid...)
	return res
}

// UnmarshalObjRef a byte slice into ObjRef
func UnmarshalObjRef(data []byte) (*ObjRef, error) {
	if len(data) != 20 {
		return nil, ErrInvalidType
	}
	res := &ObjRef{}
	copy(res.Type, data[:4])
	copy(res.Uuid, data[4:])
	return res, nil
}

// TypeId returns the ObjRef as a *TypeId
func (r *ObjRef) TypeId() *TypeId {
	res := &TypeId{}
	copy(res.data[:4], r.Type)
	copy(res.data[4:], r.Uuid)
	return res
}
